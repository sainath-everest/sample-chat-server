package security

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sainath-everest/sample-chat-server/model"

	"github.com/sainath-everest/sample-chat-server/database"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

func Signin(w http.ResponseWriter, r *http.Request) {
	var loginStatus string
	log.Println("User SignIn")
	var signedUser model.User
	err := json.NewDecoder(r.Body).Decode(&signedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var expectedUser model.User = database.GetUserByID(signedUser.UserID)

	if (model.User{}) == expectedUser || expectedUser.Password != signedUser.Password {
		//w.WriteHeader(http.StatusUnauthorized)
		//return
		loginStatus = "fail"
	} else {
		loginStatus = "success"

	}
	if loginStatus == "success" {
		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &model.Claims{
			UserID: signedUser.UserID,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

	}

	log.Println(loginStatus)
	w.Write([]byte(loginStatus))

}
func ValidateToken(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("token")
	log.Println("token", c)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	tknStr := c.Value

	claims := &model.Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.UserID)))
	return true
}
