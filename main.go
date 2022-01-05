package main

import (
	"database/sql"
	"log"
	"net/http"

	"KZ_forum/data"
	"KZ_forum/handler"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// query := url.Values{
	// 	"id": []string{},
	// }
	// fmt.Println(query.Encode())
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	data.CreateTables(db)
	testPosts := []data.Post{
		{Title: "Title1", Content: "Content1"},
		{Title: "Title2", Content: "Content2"},
		{Title: "Title3", Content: "Content3"},
		{Title: "Title4", Content: "Content4"},
		{Title: "Title5", Content: "Content5"},
		{Title: "Title6", Content: "Content6"},
		{Title: "Title7", Content: "Content7"},
	}
	data.InsertPosts(db, testPosts)

	http.HandleFunc("/", handler.IndexHandler)
	http.HandleFunc("/post", handler.PostHandler)
	http.HandleFunc("/write", handler.WriteHandler)
	http.HandleFunc("/signup", handler.SignupHandler)
	http.HandleFunc("/signup-submit", handler.SignupSubmitHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/login-submit", handler.LoginSubmitHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
	http.ListenAndServe(":8888", nil)
}
