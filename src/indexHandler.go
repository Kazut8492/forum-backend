package src

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))
	if r.URL.Path != "/" {
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
	posts := ReadPosts(db)

	if err := tpl.ExecuteTemplate(w, "index.html", posts); err != nil {
		w.WriteHeader(500)
		return
	}
}
