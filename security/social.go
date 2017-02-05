package security

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"bytes"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/twitter"
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

var Config SecurityConfig

func Configure(cfg SecurityConfig) {
	Config = cfg
	log.Println("Reading configuration\n")
	log.Println(Config.String())
	gothic.Store = sessions.NewCookieStore([]byte(Config.JWTSecret))
	goth.UseProviders(
		twitter.New(cfg.Social["twitter"].ClientID,
			cfg.Social["twitter"].Secret,
			"https://"+cfg.BaseURL+"/callback?provider=twitter"),
		github.New(cfg.Social["github"].ClientID,
			cfg.Social["github"].Secret,
			"https://"+cfg.BaseURL+"/callback?provider=github"),
	)
}

func SocialCallbackHandler(rw http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(rw, req)
	if err != nil {
		log.Println(err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":        fmt.Sprintf("%s@%s", user.UserID, user.Provider),
		"FirstName": user.FirstName,
		"LastName":  user.LastName,
		"expires":   time.Now().Add(72 * time.Hour).Unix(),
	})
	tkn, err := token.SignedString([]byte(Config.JWTSecret))
	session, _ := gothic.Store.Get(req, gothic.SessionName)
	session.Values[gothic.SessionName] = ""
	session.Options.MaxAge = -1
	session.Save(req, rw)
	http.Redirect(rw, req, "/login?token="+tkn, http.StatusPermanentRedirect)
}
