package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/like-post" {
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
	fromPage := r.FormValue("from")
	postID := r.FormValue("id")

	// Check if the user is logged-in. If cookie is empty, redirect to the index page.
	cookie, err := r.Cookie("session")
	receivedUUID := cookie.Value
	matchedUsername := getUsernameFromUUID(w, receivedUUID)
	if err != nil || matchedUsername == "" {
		fmt.Println("ERROR: Log-in needed to react to a post")
		switch {
		case fromPage == "index":
			http.Redirect(w, r, "/", http.StatusFound)
		case fromPage == "post":
			http.Redirect(w, r, "/post?id="+postID, http.StatusFound)
		}
		return
	}

	// Insert like
	statement, err := db.Prepare(`
		INSERT INTO like (
			post_id,
			creator_username
		) VALUES (?, ?)
	`)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer statement.Close()
	// number of variables have to be matched with INSERTed variables
	statement.Exec(postID, matchedUsername)

	switch {
	case fromPage == "index":
		http.Redirect(w, r, "/", http.StatusFound)
	case fromPage == "post":
		http.Redirect(w, r, "/post?id="+postID, http.StatusFound)
	}
}
