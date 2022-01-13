package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	defer db.Close()

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	user := User{}
	user.Username = username
	user.Pass = password
	// userData := ReadUser(w, db, user)

	var matchedUsername string
	db.QueryRow("SELECT username FROM user WHERE username = ?", username).Scan(&matchedUsername)
	if matchedUsername == "" {
		//WORK IN PROGRESS
		fmt.Println("ERROR: log in failed, username not found in the database")
		http.Redirect(w, r, "/login", http.StatusFound)
		// return
	} else {
		statement, err := db.Query("SELECT * FROM user WHERE username = ?", username)
		if err != nil {
			w.WriteHeader(500)
			fmt.Println(err.Error())
			log.Fatal(1)
		}
		defer statement.Close()

		user := User{}
		for statement.Next() {
			err := statement.Scan(&user.ID, &user.Username, &user.Email, &user.Pass)
			if err != nil {
				w.WriteHeader(500)
				fmt.Println(err.Error())
				log.Fatal(1)
			}
		}
		if err := CompareHashAndPassword(user.Pass, password); err != nil {
			//WORK IN PROGRESS
			fmt.Println("ERROR: log in failed, password not matched")
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			fmt.Println("log in successed")
			InitiateSession(w, r, db, user)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}
