package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sainath-everest/sample-chat-server/model"
	"github.com/sainath-everest/sample-chat-server/security"
)

var id string

func handleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	fmt.Println("i am new connection ")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ids, ok := r.URL.Query()["id"]

	token, ok := r.URL.Query()["token"]
	log.Println("debug ", token)

	if !ok {
		log.Println("Url Param 'key' is missing")
		return
	}
	id := ids[0]
	isValidToken := security.ValidateToken(w, r, token[0])
	if isValidToken {
		log.Println("token is valid")
		ws, _ := upgrader.Upgrade(w, r, nil)
		client := &Client{ID: id, Hub: hub, Conn: ws}
		client.Hub.Register <- client

		for {
			var msg model.Message
			err := client.Conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				break
			} else {
				msg.MessageType = "incoming"
			}

			client.Hub.Send <- msg
		}

	}

}
