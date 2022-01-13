package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/like-comment" {
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
	postID := r.FormValue("postid")
	commentID := r.FormValue("commentid")
	db.Exec("UPDATE comment SET like = like + 1 WHERE post_id = ? AND comment_id = ?", postID, commentID)

	http.Redirect(w, r, "/post?id="+postID, http.StatusFound)
}
