package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleConnections(hub *hub, w http.ResponseWriter, r *http.Request) {
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
	// Make sure we close the connection when the function returns
	// defer ws.Close()

	client := &client{ID: id, hub: hub, conn: ws}

	// Register our new client
	client.hub.register <- client

	for {
		var msg message
		// Read in a new message as JSON and map it to a Message object
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			// delete(hub.connections, client)
			break
		}
		// Send the newly received message to the broadcast channel
		client.hub.send <- msg
	}

}
