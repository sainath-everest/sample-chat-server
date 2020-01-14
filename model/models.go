package model

import (
	"github.com/dgrijalva/jwt-go"
)

type Message struct {
	SenderID    string `json:"senderId"`
	ReceiverID  string `json:"receiverId"`
	Data        string `json:"data"`
	Date        string `json:"date"`
	MessageType string `json:"messageType"`
}
type UserLoginStatus struct {
	Token          string `json:"token"`
	IsLoginSuccess string `json:"isLoginSuccess"`
}
type User struct {
	UserID          string
	FirstName       string
	LastName        string
	Password        string
	ConfirmPassword string
}

type Claims struct {
	UserID string `json:"username"`
	jwt.StandardClaims
}
