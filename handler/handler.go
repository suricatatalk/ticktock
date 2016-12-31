package handler

import (
	"log"
	"net/http"

	"encoding/json"

	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sohlich/ticktock/domain"
	"github.com/sohlich/ticktock/logic"
	"github.com/sohlich/ticktock/security"
)

type SecuredHandler func(user security.User, rw http.ResponseWriter, req *http.Request)

func JWTAuthHandler(h SecuredHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appendHeaders(w)
		tkn := r.Header.Get("x-auth")
		if tkn == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
			return []byte(security.Config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
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

func appendHeaders(w http.ResponseWriter) {
	w.Header().Set("access-control-expose-headers", "x-auth")
}

func Tasks(user security.User, rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getTasks(user, rw, req)
	}
}

func getTasks(user security.User, rw http.ResponseWriter, req *http.Request) {
	all, err := domain.Tasks.FindAllByOwner(user.ID)
	if err != nil {
		log.Printf("Error while loading tasks: " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}
	encoder := json.NewEncoder(rw)
	encoder.Encode(all)
}

func Events(user security.User, rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodPut {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Decode posted event
	event := &logic.EventDTO{}
	if err := json.NewDecoder(req.Body).Decode(event); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}
	defer req.Body.Close()

	var eventHandler logic.EventFunction
	switch strings.ToLower(event.EventTypeString) {
	case "start":
		eventHandler = logic.Start
	case "pause":
		eventHandler = logic.Pause
	}

	task, err := eventHandler(user, event)
	if err != nil {
		log.Printf("Error occured : %s\n", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(rw).Encode(task)
}
