package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sainath-everest/sample-chat-server/database"
	"github.com/sainath-everest/sample-chat-server/model"
)

type Hub struct {
	Clients map[string]*Client

	// Inbound messages from the clients.
	Send chan model.Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	UnsentMessageMap map[string][]model.Message
}

func newHub() *Hub {
	return &Hub{
		Send:       make(chan model.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}
}

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
	ID   string
	Hub  *Hub
}

func sendOfflineMessages(client *Client) {
	log.Println("in sendOfflineMessages", client.ID)
	messages := database.GetOfflineMessages(client.ID)
	for index, message := range messages {
		fmt.Printf("%v: %v\n", index, message)
		client.Hub.Send <- message
	}
	database.DeleteOfflineMessages(client.ID)
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.ID] = client
			go sendOfflineMessages(client)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)

			}
		case message := <-h.Send:
			log.Println(" sending  message to  ", message.ReceiverID)
			if client, ok := h.Clients[message.ReceiverID]; ok {
				log.Println("before write message to client")
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Println(err)
					log.Println("offline test before adding offline msg to map ", message)
					database.StoreOfflineMessages(message)
					client.Conn.Close()
				}

			}
		}
	}
}
