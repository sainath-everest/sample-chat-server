package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func handleUserRegistration(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:sai@test@tcp(127.0.0.1:3306)/testDb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var newUser user
	reqBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &newUser)
	log.Println("newUser ", newUser)
	insForm, err := db.Prepare("INSERT INTO user(user_id,first_name,last_name,pass_word,creation_date) VALUES(?,?,?,?,NOW())")
	insForm.Exec(newUser.UserID, newUser.FirstName, newUser.LastName, newUser.Password)

	if err != nil {
		panic(err.Error())
	}

}
func getUserByID(userID string) user {
	db, err := sql.Open("mysql", "root:sai@test@tcp(127.0.0.1:3306)/testDb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var user user
	err = db.QueryRow("SELECT user_id,pass_word FROM user where user_id= ?", userID).Scan(&user.UserID, &user.Password)
	return user

}
