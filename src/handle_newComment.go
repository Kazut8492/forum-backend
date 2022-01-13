package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/new-comment" {
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
	postIDstr := r.FormValue("id")
	postID, err := strconv.Atoi(postIDstr)
	// fmt.Println(postID)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Check if the user is logged-in. If cookie is empty, redirect to the post page.
	cookie, err := r.Cookie("session")
	receivedUUID := cookie.Value
	matchedUsername := getUsernameFromUUID(w, receivedUUID)
	if err != nil || receivedUUID != matchedUsername {
		fmt.Println("ERROR: Log-in needed to create a comment")
		http.Redirect(w, r, "/post?id="+postIDstr, http.StatusFound)
		return
	}

	var newComment Comment
	newComment.PostId = postID
	newComment.Title = r.FormValue("commentTitle")
	newComment.Content = r.FormValue("commentDescription")
	newComment.CreatorUsrName = matchedUsername
	InsertComment(db, newComment)
	http.Redirect(w, r, "/post?id="+postIDstr, http.StatusFound)
}
