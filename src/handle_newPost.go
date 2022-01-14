package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/new-post" {
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

	// Check if the user is logged-in. If cookie is empty, redirect to the index page.
	cookie, err := r.Cookie("session")
	receivedUUID := cookie.Value
	matchedUsername := getUsernameFromUUID(w, receivedUUID)
	if err != nil || matchedUsername == "" {
		fmt.Println("ERROR: Log-in needed to create a post")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	r.ParseForm()
	postTitle := r.FormValue("postTitle")
	postContent := r.FormValue("postContent")
	postCategory := r.Form["category"]
	// At least one category has to be selected. Otherwise, redirecting to the index page.
	if len(postCategory) == 0 {
		fmt.Println("ERROR: At least one category has to be selected")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var post Post
	post.Title = postTitle
	post.Content = postContent
	post.CategoryArr = postCategory
	post.CreatorUsrName = matchedUsername
	InsertPost(db, post)
	http.Redirect(w, r, "/", http.StatusFound)
}
