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
		loginStatus = "fail"
	} else {
		loginStatus = "success"

	}
	if loginStatus == "success" {
		expirationTime := time.Now().Add(2160 * time.Hour)
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
		log.Println(tokenString)

		userLoginStatus := model.UserLoginStatus{tokenString, loginStatus}

		js, err := json.Marshal(userLoginStatus)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}

}
func ValidateToken(w http.ResponseWriter, r *http.Request, token string) bool {

	tknStr := token
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
	return true
}
