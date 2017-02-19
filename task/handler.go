package task

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/sohlich/ticktock/user"
)

// Handler to dispatch task requests
func TasksHandler(user user.User, rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		getTasksHandler(user, rw, req)
	}
}

func getTasksHandler(user user.User, rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	limit, err := strconv.ParseInt(req.Form.Get("limit"), 10, 32)
	if err != nil {
		limit = 10
	}
	all, err := Repository.FindAllByOwner(user.ID, int(limit))
	if err != nil {
		log.Printf("Error while loading tasks: " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(rw)
	encoder.Encode(all)
}

// Events handler provides api for events
func EventsHandler(user user.User, rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodPut {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Decode posted event
	event := &Event{}
	if err := json.NewDecoder(req.Body).Decode(event); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	log.Printf("Handling event: %v\n", event)

	var eventHandler EventFunction
	switch strings.ToLower(event.EventType) {
	case "start":
		if len(event.TaskID) == 0 {
			eventHandler = Start
		} else {
			eventHandler = Resume
		}
	case "pause":
		eventHandler = Pause
	case "finish":
		eventHandler = Finish
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

// Handles the tags
func TagsHandler(user user.User, rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodPut {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	switch req.Method {
	case http.MethodPost:
		t := &Task{}
		if err := json.NewDecoder(req.Body).Decode(t); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()
		Repository.InsertTags(t.ID, t.Tags)
	}
}
