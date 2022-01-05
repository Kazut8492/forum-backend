package handler

import (
	"KZ_forum/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func SignupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	userPass := r.FormValue("userPass")
	encryptedUserPass, err := data.PasswordEncrypt(userPass)
	if err != nil {
		fmt.Println("ERROR: Failed to encrypt passworrd")
		log.Fatal(1)
	}
	newUser := data.User{}
	newUser.Username = username
	newUser.Pass = encryptedUserPass
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	data.InsertUser(db, newUser)
	http.Redirect(w, r, "/", http.StatusFound)
}
