package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sainath-everest/sample-chat-server/database"
	"github.com/sainath-everest/sample-chat-server/security"

	"github.com/gorilla/websocket"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi there!")
	})
	http.HandleFunc("/registration", func(w http.ResponseWriter, r *http.Request) {
		database.HandleUserRegistration(w, r)
	})
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		security.Signin(w, r)
	})
	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		security.Welcome(w, r)
	})

	// Configure websocket route
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		//enableCors(&w)
		handleConnections(hub, w, r)

	})

	// Start listening for incoming chat messages
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// err := httpscerts.Check("cert.pem", "key.pem")
	// if err != nil {
	// 	err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:8081")
	// 	if err != nil {
	// 		log.Fatal("Error: Couldn't create https certs.")
	// 	}
	// }
	// http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", nil)

}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
