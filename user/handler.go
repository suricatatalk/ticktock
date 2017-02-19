package user

import (
	"encoding/json"
	"log"
	"net/http"
)

func UserHandler(u User, w http.ResponseWriter, r *http.Request) {
	log.Println("UserHandler")
	switch r.Method {
	case http.MethodGet:
		if b, err := json.Marshal(u); err == nil {
			log.Println("Serving: " + string(b))
			w.Write(b)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Err: " + err.Error())
		}
	case http.MethodPost:
		form := &User{}
		d := json.NewDecoder(r.Body)
		if err := d.Decode(form); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Saving: ", form)
		if err := Repository.Save(form); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if b, err := json.Marshal(form); err == nil {
			w.Write(b)
		} else {
			log.Println("Err: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
