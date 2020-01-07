package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sainath-everest/sample-chat-server/database"
	"github.com/sainath-everest/sample-chat-server/security"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi there!")
	})
	http.HandleFunc("/registration", func(w http.ResponseWriter, r *http.Request) {
		database.HandleUserRegistration(w, r)
	})
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		security.Signin(w, r)
	})
	http.HandleFunc("/getAllUsers", func(w http.ResponseWriter, r *http.Request) {
		database.GetAllUsers(w, r)
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, w, r)

	})

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
