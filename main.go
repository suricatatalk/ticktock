package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"encoding/json"

	"strings"

	"github.com/braintree/manners"
	"github.com/gorilla/context"
	"github.com/markbates/goth/gothic"
	"github.com/sohlich/ticktock/domain"
	"github.com/sohlich/ticktock/handler"
	"github.com/sohlich/ticktock/security"
)

func main() {
	configureApp()

	mux := http.NewServeMux()
	appHandler := http.FileServer(http.Dir("/Users/radek/Projects/Html/ticktock/ticktock/dist"))
	mux.HandleFunc("/auth", gothic.BeginAuthHandler)
	mux.HandleFunc("/callback", security.SocialCallbackHandler)
	mux.HandleFunc("/tasks", handler.JWTAuthHandler(handler.Tasks))
	mux.HandleFunc("/events", handler.JWTAuthHandler(handler.Events))
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

func configureApp() {
	config := security.SecurityConfig{}
	domainCfg := domain.StorageConfig{}
	f, err := os.Open("config.json")
	if err != nil {
		log.Println(err.Error())
	}
	enc := json.NewDecoder(f)
	enc.Decode(&config)
	enc.Decode(&domainCfg)
	f.Close()

	security.Configure(config)
	domain.Open(domainCfg)
}
