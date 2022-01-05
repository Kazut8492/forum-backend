package handler

import (
	"KZ_forum/data"
	"database/sql"
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	posts := data.ReadPosts(db)

	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}
	if err := tpl.ExecuteTemplate(w, "index.html", posts); err != nil {
		w.WriteHeader(500)
		return
	}
}
