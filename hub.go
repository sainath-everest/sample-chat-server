package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type message struct {
	ID   string
	Data string
}

type client struct {
	conn *websocket.Conn
	mu   sync.Mutex
	ID   string
	hub  *hub
} 

type hub struct {
	clients map[string]*client

	// Inbound messages from the clients.
	send chan message

	// Register requests from the clients.
	register chan *client

	// Unregister requests from clients.
	unregister chan *client

	connections map[string]*client
}

func newHub() *hub {
	return &hub{
		send:       make(chan message),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[string]*client),
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)

			}
		case message := <-h.send:
			if client, ok := h.clients[message.ID]; ok {
				fmt.Println(message, "hihow")
				err := client.conn.WriteJSON(message)
				if err != nil {
					log.Printf("error: %v", err)
					client.conn.Close()
				}
			}
		}
	}
}
