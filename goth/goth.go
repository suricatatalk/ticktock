package goth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"bytes"

	"net/url"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/twitter"
	"github.com/sohlich/ticktock/user"
)

type OAuthConfig struct {
	ClientID string
	Secret   string
}

type SecurityConfig struct {
	Social    map[string]OAuthConfig
	JWTSecret string
	BaseURL   string
}

func (c *SecurityConfig) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("Configured social:\n")
	for prov, _ := range c.Social {
		buf.WriteString(prov + "\n")
	}
	return string(buf.Bytes())
}

var cfg SecurityConfig

func InitGoth() {
	cfg := SecurityConfig{}
	readFileToStruct("config.json", &cfg)
	log.Println(cfg.String())
	gothic.Store = sessions.NewCookieStore([]byte(cfg.JWTSecret))
	goth.UseProviders(
		twitter.New(cfg.Social["twitter"].ClientID, cfg.Social["twitter"].Secret, "https://"+cfg.BaseURL+"/callback?provider=twitter"),
		github.New(cfg.Social["github"].ClientID, cfg.Social["github"].Secret, "https://"+cfg.BaseURL+"/callback?provider=github"),
	)
}

func readFileToStruct(file string, cfg interface{}) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err.Error())
	}
	enc := json.NewDecoder(bytes.NewReader(b))
	enc.Decode(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func SocialCallbackHandler(rw http.ResponseWriter, req *http.Request) {
	u, err := gothic.CompleteUserAuth(rw, req)
	if err != nil || len(u.UserID) == 0 {
		log.Println(err)
		http.Redirect(rw, req, "/login?error="+url.QueryEscape(err.Error()), http.StatusPermanentRedirect)
		return
	}

	ID := fmt.Sprintf("%s@%s", u.UserID, u.Provider)

	// try to fill
	var dUser *user.User
	dUser, err = user.Repository.FindById(ID)
	var f string
	if len(dUser.ID) == 0 {
		f = "&firstlogin=true"
	} else {
		u.FirstName = dUser.Firstname
		u.LastName = dUser.Lastname
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":        ID,
		"FirstName": u.FirstName,
		"LastName":  u.LastName,
		"expires":   time.Now().Add(72 * time.Hour).Unix(),
	})
	tkn, err := token.SignedString([]byte(cfg.JWTSecret))
	session, _ := gothic.Store.Get(req, gothic.SessionName)
	session.Values[gothic.SessionName] = ""
	session.Options.MaxAge = -1
	session.Save(req, rw)

	http.Redirect(rw, req, "/login?token="+tkn+f, http.StatusPermanentRedirect)
}

type SecuredHandler func(user user.User, rw http.ResponseWriter, req *http.Request)

func JWTAuthHandler(h SecuredHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appendHeaders(w)
		tkn := r.Header.Get("x-auth")
		if tkn == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user := user.User{}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			user.ID, _ = claims["ID"].(string)
			user.Firstname, _ = claims["Firstname"].(string)
			user.Lastname, _ = claims["LastName"].(string)
		}

		h(user, w, r)
	}
}

func appendHeaders(w http.ResponseWriter) {
	w.Header().Set("access-control-expose-headers", "x-auth")
}
