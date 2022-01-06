package src

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Jump to a certain post selected on the index page
func PostHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	posts := ReadPosts(db)

	postIDstr := r.FormValue("id")
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		log.Fatal(err.Error())
	}

	comments := ReadComments(db, postID)

	certainPost := Post{
		ID:       posts[postID-1].ID,
		Title:    posts[postID-1].Title,
		Content:  posts[postID-1].Content,
		Comments: comments,
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
