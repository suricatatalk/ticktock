package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"strings"

	"github.com/braintree/manners"
	"github.com/gorilla/context"
	"github.com/markbates/goth/gothic"
	"github.com/sohlich/ticktock/config"
	"github.com/sohlich/ticktock/goth"
	"github.com/sohlich/ticktock/task"
	"github.com/sohlich/ticktock/user"
)

func main() {
	log.Printf("Starting application %s", "TickTock")
	config.InitDB("database.json", "Development")
	goth.InitGoth()

	task.InitializeRepository(config.DB)
	user.InitializeRepository(config.DB)

	mux := http.NewServeMux()
	appHandler := http.FileServer(http.Dir("/Users/radek/Projects/Html/ticktock/ticktock/dist"))
	mux.HandleFunc("/auth", gothic.BeginAuthHandler)
	mux.HandleFunc("/callback", goth.SocialCallbackHandler)
	mux.HandleFunc("/tasks", goth.JWTAuthHandler(task.TasksHandler))
	mux.HandleFunc("/events", goth.JWTAuthHandler(task.EventsHandler))
	mux.HandleFunc("/tags", goth.JWTAuthHandler(task.TagsHandler))
	mux.HandleFunc("/user", goth.JWTAuthHandler(user.UserHandler))
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("Handling %v", req.URL)
		if !strings.Contains(req.URL.Path, ".") {
			req.URL.Path = "/"
		}
		appHandler.ServeHTTP(rw, req)
	})

	server := manners.NewServer()
	server.Handler = context.ClearHandler(mux)

	errChan := make(chan (error))
	go func() {
		server.Addr = ":7070"
		errChan <- server.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			server.BlockingClose()
			os.Exit(0)
		}
	}
}
