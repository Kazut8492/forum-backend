package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func SignupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup-submit" {
		w.WriteHeader(404)
		return
	}
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	defer db.Close()

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
		fmt.Println("ERROR: Same email or username found in the database")
		http.Redirect(w, r, "/signup", http.StatusFound)
	} else {
		user := User{}
		user.Username = username
		user.Email = email
		user.Pass = encryptedPass
		InsertUser(db, user)
		InitiateSession(w, r, db, user)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
