package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/registration", func(w http.ResponseWriter, r *http.Request) {
		handleUserRegistration(w, r)
	})
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		signin(w, r)
	})
	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		welcome(w, r)
	})

	// Configure websocket route
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(hub, w, r)
	})

	// Start listening for incoming chat messages
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
