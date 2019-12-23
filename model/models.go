package model

import (
	"github.com/dgrijalva/jwt-go"
)

type Message struct {
	ID   string
	Data string
}
type User struct {
	UserID    string
	FirstName string
	LastName  string
	Password  string
}

type Claims struct {
	UserID string `json:"username"`
	jwt.StandardClaims
}
