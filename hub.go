package main

import (
	"log"
)

func newHub() *hub {
	return &hub{
		send:             make(chan message),
		register:         make(chan *client),
		unregister:       make(chan *client),
		clients:          make(map[string]*client),
		unsentMessageMap: make(map[string][]message),
	}
}

func sendOfflineMessages(client *client) {
	for k, v := range client.hub.unsentMessageMap {
		if k == client.ID {
			for i := 0; i < len(v); i++ {
				client.hub.send <- v[i]

			}
			delete(client.hub.unsentMessageMap, client.ID)
		}
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
			go sendOfflineMessages(client)

		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)

			}
		case message := <-h.send:
			log.Println("message id ", message.ID)
			if client, ok := h.clients[message.ID]; ok {
				err := client.conn.WriteJSON(message)
				if err != nil {
					h.unsentMessageMap[message.ID] = append(h.unsentMessageMap[message.ID], message)
					client.conn.Close()
				}

			}
		}
	}
}
