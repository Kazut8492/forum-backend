package handler

import (
	"KZ_forum/data"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	posts := data.ReadPosts(db)

	postIDstr := r.FormValue("id")
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		log.Fatal(err.Error())
	}

	comments := data.ReadComments(db, postID)

	certainPost := data.Post{posts[postID-1].Id, posts[postID-1].Title, posts[postID-1].Content, comments}

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
