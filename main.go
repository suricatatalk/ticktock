package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"strings"

	"github.com/braintree/manners"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/twitter"
	"github.com/sohlich/ticktock/handler"
	"github.com/sohlich/ticktock/model"
)

func main() {
	log.Printf("Starting application %s", "TickTock")
	model.InitDB("database.json", "Development")
	handler.InitGoth()

	mux := http.NewServeMux()
	appHandler := http.FileServer(http.Dir("/Users/radek/Projects/Html/ticktock/ticktock/dist"))
	mux.HandleFunc("/auth", gothic.BeginAuthHandler)
	mux.HandleFunc("/callback", handler.SocialCallbackHandler)
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

func InitDB() {

}

func InitGoth() {
	cfg := handler.SecurityConfig{}
	readFileToStruct("config.json", &cfg)
	log.Println(cfg.String())
	gothic.Store = sessions.NewCookieStore([]byte(cfg.JWTSecret))
	goth.UseProviders(
		twitter.New(cfg.Social["twitter"].ClientID, cfg.Social["twitter"].Secret, "https://"+cfg.BaseURL+"/callback?provider=twitter"),
		github.New(cfg.Social["github"].ClientID, cfg.Social["github"].Secret, "https://"+cfg.BaseURL+"/callback?provider=github"),
	)
}

func readFileToStruct(file string, cfg interface{}) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err.Error())
	}
	enc := json.NewDecoder(bytes.NewReader(b))
	enc.Decode(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
