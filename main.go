package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"encoding/json"

	"github.com/braintree/manners"
	"github.com/gorilla/context"
	"github.com/markbates/goth/gothic"
	"github.com/sohlich/ticktock/handler"
	"github.com/sohlich/ticktock/security"
)

func main() {
	configureApp()
	mux := http.NewServeMux()
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/callback", security.SocialCallbackHandler)
	mux.HandleFunc("/auth", gothic.BeginAuthHandler)

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

func configureApp() {
	config := security.SecurityConfig{}
	f, err := os.Open("config.json")
	if err != nil {
		log.Println(err.Error())
	}
	enc := json.NewDecoder(f)
	enc.Decode(&config)
	f.Close()

	security.Configure(config)
}
