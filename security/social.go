package security

import (
	"fmt"
	"log"
	"net/http"
	"time"

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
}

var config SecurityConfig

func init() {
	gothic.Store = sessions.NewCookieStore([]byte("goth-example"))

}

func Configure(cfg SecurityConfig) {
	config = cfg
	goth.UseProviders(
		twitter.New(cfg.Social["twitter"].ClientID,
			cfg.Social["twitter"].Secret,
			"https://71d90784.ngrok.io/callback?provider=twitter"),
		github.New(cfg.Social["github"].ClientID,
			cfg.Social["github"].Secret,
			"https://71d90784.ngrok.io/callback?provider=github"),
	)
}

func SocialCallbackHandler(rw http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(rw, req)
	if err != nil {
		log.Println(err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":  fmt.Sprintf("%s@%s", user.UserID, user.Provider),
		"expires": time.Now().Add(72 * time.Hour).Unix(),
	})
	tkn, err := token.SignedString([]byte("secret"))
	rw.Header().Set("x-auth", tkn)
}
