package src

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))
	if r.URL.Path != "/filter" {
		w.WriteHeader(404)
		return
	}
	r.ParseForm()
	appliedFilter := r.Form["filter"]
	// At least one category has to be selected. Otherwise, redirecting to the index page.
	if len(appliedFilter) == 0 {
		fmt.Println("ERROR: At least one category has to be selected")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	defer db.Close()
	fullPosts := ReadPosts(db)

	filteredPosts := []Post{}
	for index, post := range fullPosts {
		for _, postCategory := range post.CategoryArr {
			if contains(appliedFilter, postCategory) {
				filteredPosts = append(filteredPosts, fullPosts[index])
				break
			}
		}
	}

	if err := tpl.ExecuteTemplate(w, "index.html", filteredPosts); err != nil {
		w.WriteHeader(500)
		return
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
