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
	"github.com/sohlich/ticktock/task"
)

func main() {
	log.Printf("Starting application %s", "TickTock")
	config.InitDB("database.json", "Development")
	config.InitGoth()

	task.InitializeRepository(config.DB)

	mux := http.NewServeMux()
	appHandler := http.FileServer(http.Dir("/Users/radek/Projects/Html/ticktock/ticktock/dist"))
	mux.HandleFunc("/auth", gothic.BeginAuthHandler)
	mux.HandleFunc("/callback", config.SocialCallbackHandler)
	mux.HandleFunc("/tasks", config.JWTAuthHandler(task.TasksHandler))
	mux.HandleFunc("/events", config.JWTAuthHandler(task.EventsHandler))
	mux.HandleFunc("/tags", config.JWTAuthHandler(task.TagsHandler))
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
