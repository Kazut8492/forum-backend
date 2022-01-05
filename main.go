package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"KZ_forum/src"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// Create dummy data
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	src.CreateTables(db)
	testPosts := []src.Post{
		{Title: "Title1", Content: "Content1"},
		{Title: "Title2", Content: "Content2"},
		{Title: "Title3", Content: "Content3"},
		{Title: "Title4", Content: "Content4"},
		{Title: "Title5", Content: "Content5"},
		{Title: "Title6", Content: "Content6"},
		{Title: "Title7", Content: "Content7"},
	}
	src.InsertPosts(db, testPosts)

	// Server
	http.HandleFunc("/", src.IndexHandler)
	http.HandleFunc("/post", src.PostHandler)
	http.HandleFunc("/write", src.WriteHandler)
	http.HandleFunc("/signup", src.SignupHandler)
	http.HandleFunc("/signup-submit", src.SignupSubmitHandler)
	http.HandleFunc("/login", src.LoginHandler)
	http.HandleFunc("/login-submit", src.LoginSubmitHandler)
	http.HandleFunc("/logout", src.LogoutHandler)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
	http.ListenAndServe(":8888", nil)
}
