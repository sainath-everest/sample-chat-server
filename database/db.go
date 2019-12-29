package database

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/sainath-everest/sample-chat-server/model"

	//need for mysql connection
	_ "github.com/go-sql-driver/mysql"
)

func HandleUserRegistration(w http.ResponseWriter, r *http.Request) {
	var registrationStatus string
	db, err := sql.Open("mysql", "root:sai@test@tcp(127.0.0.1:3306)/testDb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var newUser model.User
	reqBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newUser)
	log.Println("newUser ", newUser)
	var user = GetUserByID(newUser.UserID)
	if (model.User{}) == user {
		insForm, err := db.Prepare("INSERT INTO user(user_id,first_name,last_name,pass_word,creation_date) VALUES(?,?,?,?,NOW())")
		insForm.Exec(newUser.UserID, newUser.FirstName, newUser.LastName, newUser.Password)
		if err != nil {
			log.Printf(err.Error())

		}
		registrationStatus = "success"

	} else {
		registrationStatus = "fail"
	}
	log.Println(registrationStatus)
	w.Write([]byte(registrationStatus))

}
func GetUserByID(userID string) model.User {
	db, err := sql.Open("mysql", "root:sai@test@tcp(127.0.0.1:3306)/testDb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var user model.User
	err = db.QueryRow("SELECT user_id,pass_word FROM user where user_id= ?", userID).Scan(&user.UserID, &user.Password)
	return user

}
