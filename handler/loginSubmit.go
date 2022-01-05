package handler

import (
	"KZ_forum/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	userPass := r.FormValue("userPass")
	loginUser := data.User{}
	loginUser.Username = username
	loginUser.Pass = userPass
	// fmt.Printf("Input login user info: %v", loginUser)
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	userData := data.ReadUser(db, loginUser)
	if (userData == data.User{}) {
		fmt.Println("log in failed, empty data returned")
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		if err := data.CompareHashAndPassword(userData.Pass, userPass); err != nil {
			fmt.Println("log in failed, password not matched")
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			fmt.Println("log in successed")
			data.InitiateSession(w, r, db, userData)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}
