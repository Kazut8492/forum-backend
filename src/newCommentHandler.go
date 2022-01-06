package src

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	postIDstr := r.FormValue("id")
	postID, err := strconv.Atoi(postIDstr)
	// fmt.Println(postID)
	if err != nil {
		log.Fatal(err.Error())
	}

	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// defer db.Close()
	// comments := readComments(db, postID)

	var newComment Comment
	// newComment.Id = 0
	newComment.PostId = postID
	newComment.Title = r.Form["commentTitle"][0]
	newComment.Content = r.Form["commentDescription"][0]
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	InsertComments(db, newComment)
	// fmt.Println(newComment)
	// testPosts[postID-1].Comments = append(testPosts[postID-1].Comments, &newComment)
	http.Redirect(w, r, "/post?id="+postIDstr, http.StatusFound)
}
