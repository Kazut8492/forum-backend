package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func SignupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	defer db.Close()

	if r.URL.Path != "/signup-submit" {
		w.WriteHeader(404)
		return
	}

	r.ParseForm()
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	encryptedPass, err := PasswordEncrypt(password)
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		log.Fatal(1)
	}

	var matchedEmail string
	var matchedUsername string
	db.QueryRow("SELECT user_email FROM user WHERE user_email = ?", email).Scan(&matchedEmail)
	db.QueryRow("SELECT username FROM user WHERE username = ?", username).Scan(&matchedUsername)
	if email == matchedEmail || username == matchedUsername {
		//WORK IN PROGRESS
		fmt.Println("Same email or username found in the database")
		http.Redirect(w, r, "/signup", http.StatusFound)
	} else {
		user := User{}
		user.Username = username
		user.Email = email
		user.Pass = encryptedPass
		InsertUser(db, user)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
