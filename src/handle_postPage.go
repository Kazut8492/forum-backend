package src

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Jump to a certain post selected on the index page
func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	posts := ReadPosts(db)

	postIDstr := r.FormValue("id")
	// postIDstr := r.URL.Query().Get("id") this could work ????
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		log.Fatal(err.Error())
	}

	comments := ReadComments(db, postID)

	likeRows, err := db.Query(`
		SELECT * FROM like WHERE post_id = ?
	`, postID)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer likeRows.Close()
	var likes []Like
	for likeRows.Next() {
		var like Like
		err = likeRows.Scan(&like.ID, &like.PostId, &like.CommentId, &like.CreatorUsrName)
		if err != nil {
			panic(err.Error())
		}
		likes = append(likes, like)
	}

	dislikeRows, err := db.Query(`
		SELECT * FROM dislike WHERE post_id = ?
	`, postID)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer dislikeRows.Close()
	var dislikes []Dislike
	for dislikeRows.Next() {
		var dislike Dislike
		err = dislikeRows.Scan(&dislike.ID, &dislike.PostId, &dislike.CommentId, &dislike.CreatorUsrName)
		if err != nil {
			panic(err.Error())
		}
		dislikes = append(dislikes, dislike)
	}

	certainPost := Post{
		ID:          posts[postID-1].ID,
		Title:       posts[postID-1].Title,
		Content:     posts[postID-1].Content,
		CategoryArr: posts[postID-1].CategoryArr,
		Comments:    comments,
		Likes:       likes,
		Dislikes:    dislikes,
	}

	if r.URL.Path != "/post" {
		w.WriteHeader(404)
		return
	}
	if 1 <= postID && postID <= len(posts) {
		if err := tpl.ExecuteTemplate(w, "post.html", certainPost); err != nil {
			w.WriteHeader(500)
			return
		}
	} else {
		w.WriteHeader(400)
		return
	}
}
