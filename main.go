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
		{Title: "Title1", Content: "Content1", CategoryArr: []string{"science", "education"}, CreatorUsrName: "DummyUser", Like: 1, DisLike: 1},
		{Title: "Title2", Content: "Content2", CategoryArr: []string{"education", "sports"}, CreatorUsrName: "DummyUser", Like: 5, DisLike: 10},
		{Title: "Title3", Content: "Content3", CategoryArr: []string{"sports", "lifehacks"}, CreatorUsrName: "DummyUser", Like: 10, DisLike: 5},
	}

	for _, post := range testPosts {
		src.InsertPost(db, post)
	}

	// Server
	http.HandleFunc("/", src.IndexHandler)
	http.HandleFunc("/filter", src.FilterHandler)
	http.HandleFunc("/post", src.PostPageHandler)
	http.HandleFunc("/new-post", src.NewPostHandler)
	http.HandleFunc("/new-comment", src.NewCommentHandler)
	http.HandleFunc("/like-post", src.LikePostHandler)
	http.HandleFunc("/dislike-post", src.DisLikePostHandler)
	http.HandleFunc("/like-comment", src.LikeCommentHandler)
	http.HandleFunc("/dislike-comment", src.DisLikeCommentHandler)
	http.HandleFunc("/signup", src.SignupHandler)
	http.HandleFunc("/signup-submit", src.SignupSubmitHandler)
	http.HandleFunc("/login", src.LoginHandler)
	http.HandleFunc("/login-submit", src.LoginSubmitHandler)
	http.HandleFunc("/logout", src.LogoutHandler)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
	http.ListenAndServe(":8888", nil)
}
