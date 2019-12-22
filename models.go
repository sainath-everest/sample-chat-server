package main

import (
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

type message struct {
	ID   string
	Data string
}
type user struct {
	UserID    string
	FirstName string
	LastName  string
	Password  string
}

type claims struct {
	UserID string `json:"username"`
	jwt.StandardClaims
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

	unsentMessageMap map[string][]message
}
