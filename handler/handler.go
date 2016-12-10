package handler

import (
	"log"
	"net/http"

	"text/template"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sohlich/ticktock/security"
)

type SecuredHandler func(user security.User, rw http.ResponseWriter, req *http.Request)

func Login(rw http.ResponseWriter, req *http.Request) {
	//twitterLink := "<a href=\"auth?provider=twitter\">Twitter</ a>"
	//githubLink := "<a href=\"auth?provider=github\">Github</ a>"
	tmp := template.New("login.html")
	tmp.ParseFiles("static/login.html", "static/navbar.html")
	if err := tmp.Execute(rw, struct{}{}); err != nil {
		log.Println(err.Error())
	}
}

func JWTAuthHandler(h SecuredHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tkn := r.Header.Get("x-auth")
		if tkn == "" {
			redirectToLogin(w, r)
			return
		}

		token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
			return []byte(security.Config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			redirectToLogin(w, r)
			return
		}

		user := security.User{}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			user.ID, _ = claims["ID"].(string)
			user.Firstname, _ = claims["Firstname"].(string)
			user.Lastname, _ = claims["LastName"].(string)
		}

		h(user, w, r)
	}
}

func Stop(user security.User, rw http.ResponseWriter, req *http.Request) {

}

func Start(user security.User, rw http.ResponseWriter, req *http.Request) {

}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", 302)
}
