package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func DisLikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dislike-post" {
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
	postID := r.FormValue("id")
	db.Exec("UPDATE post SET dislike = dislike + 1 WHERE post_id = ?", postID)

	fromPage := r.FormValue("from")
	switch {
	case fromPage == "index":
		http.Redirect(w, r, "/", http.StatusFound)
	case fromPage == "post":
		http.Redirect(w, r, "/post?id="+postID, http.StatusFound)
	}
}
