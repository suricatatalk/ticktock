package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

	base := "/Users/radek/Projects/Html/ticktock/ticktock/dist"

	mux := http.NewServeMux()
	mux.HandleFunc("/auth", gothic.BeginAuthHandler)
	mux.HandleFunc("/callback", goth.SocialCallbackHandler)
	mux.HandleFunc("/tasks", goth.JWTAuthHandler(task.TasksHandler))
	mux.HandleFunc("/events", goth.JWTAuthHandler(task.EventsHandler))
	mux.HandleFunc("/tags", goth.JWTAuthHandler(task.TagsHandler))
	mux.HandleFunc("/user", goth.JWTAuthHandler(user.UserHandler))
	mux.Handle("/", &config.WebApp{base})

	srv := &http.Server{Addr: ":7070", Handler: mux}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)

	log.Println("Server gracefully stopped")
}
