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

	r.ParseForm()
	postTitle := r.FormValue("postTitle")
	postContent := r.FormValue("postContent")
	postCategory := r.Form["category"]
	// At least one category has to be selected
	if len(postCategory) == 0 {
		fmt.Println("ERROR: At least one category has to be selected")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	var post Post
	post.Title = postTitle
	post.Content = postContent
	post.CategoryArr = postCategory
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	InsertPost(db, post)
	http.Redirect(w, r, "/", http.StatusFound)
}
