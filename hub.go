package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
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
		Send:             make(chan model.Message),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		Clients:          make(map[string]*Client),
		UnsentMessageMap: make(map[string][]model.Message),
	}
}

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
	ID   string
	Hub  *Hub
}

func sendOfflineMessages(client *Client) {
	for k, v := range client.Hub.UnsentMessageMap {
		if k == client.ID {
			for i := 0; i < len(v); i++ {
				client.Hub.Send <- v[i]

			}
			delete(client.Hub.UnsentMessageMap, client.ID)
		}
	}
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
			log.Println("message id ", message.ID)
			if client, ok := h.Clients[message.ID]; ok {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					h.UnsentMessageMap[message.ID] = append(h.UnsentMessageMap[message.ID], message)
					client.Conn.Close()
				}

			}
		}
	}
}
