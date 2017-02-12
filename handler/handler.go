package handler

import (
	"log"
	"net/http"

	"encoding/json"

	"strings"

	"strconv"

	"github.com/sohlich/ticktock/logic"
	"github.com/sohlich/ticktock/model"
)

// Handler to dispatch task requests
func Tasks(user model.User, rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getTasks(user, rw, req)
	}
}

func getTasks(user model.User, rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	limit, err := strconv.ParseInt(req.Form.Get("limit"), 10, 32)
	if err != nil {
		limit = 10
	}
	all, err := model.Tasks.FindAllByOwner(user.ID, int(limit))
	if err != nil {
		log.Printf("Error while loading tasks: " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}
	encoder := json.NewEncoder(rw)
	encoder.Encode(all)
}

// Events handler provides api for events
func Events(user model.User, rw http.ResponseWriter, req *http.Request) {
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

	log.Printf("Handling event: %v\n", event)

	var eventHandler logic.EventFunction
	switch strings.ToLower(event.EventTypeString) {
	case "start":
		if len(event.TaskID) == 0 {
			eventHandler = logic.Start
		} else {
			eventHandler = logic.Resume
		}
	case "pause":
		eventHandler = logic.Pause
	case "finish":
		eventHandler = logic.Finish
	default:
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := eventHandler(user, event)
	if err != nil {
		log.Printf("Error occured : %s\n", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(rw).Encode(task)
}
