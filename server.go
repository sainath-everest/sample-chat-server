package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sainath-everest/sample-chat-server/model"
)

var id string

func handleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	fmt.Println("i am new connection ")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ids, ok := r.URL.Query()["id"]

	if !ok {
		log.Println("Url Param 'key' is missing")
		return
	}
	id := ids[0]
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &Client{ID: id, Hub: hub, Conn: ws}

	// Register our new client
	client.Hub.Register <- client

	for {
		var msg model.Message
		// Read in a new message as JSON and map it to a Message object
		err := client.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		// Send the newly received message to the receiver channel
		client.Hub.Send <- msg
	}

}
