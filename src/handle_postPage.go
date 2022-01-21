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
		w.WriteHeader(500)
		log.Fatal(err.Error())
	}
	defer db.Close()
	posts := ReadPosts(db)

	postIDstr := r.FormValue("id")
	// postIDstr := r.URL.Query().Get("id") this could work ????
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		w.WriteHeader(500)
		log.Fatal(err.Error())
	}

	certainPost := Post{
		ID:          posts[postID-1].ID,
		Title:       posts[postID-1].Title,
		Content:     posts[postID-1].Content,
		CategoryArr: posts[postID-1].CategoryArr,
		Comments:    posts[postID-1].Comments,
		Likes:       posts[postID-1].Likes,
		Dislikes:    posts[postID-1].Dislikes,
	}

	if r.URL.Path != "/post" {
		w.WriteHeader(404)
		return
	}

	if 1 <= postID && postID <= len(posts) {
		cookie, err := r.Cookie("session")
		if err != nil {
			if err := tpl.ExecuteTemplate(w, "post.html", certainPost); err != nil {
				w.WriteHeader(500)
				return
			}
			return
		}
		receivedUUID := cookie.Value
		matchedUsername := getUsernameFromUUID(w, receivedUUID)
		if matchedUsername == "" {
			if err := tpl.ExecuteTemplate(w, "post.html", certainPost); err != nil {
				w.WriteHeader(500)
				return
			}
			return
		}

		if err := tpl.ExecuteTemplate(w, "logged-post.html", certainPost); err != nil {
			w.WriteHeader(500)
			return
		}
	} else {
		w.WriteHeader(400)
		return
	}
}
