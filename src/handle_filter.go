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
	appliedFilter := r.FormValue("filter")

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	defer db.Close()
	fullPosts := ReadPosts(db)
	filteredPosts := []Post{}

	switch {
	case appliedFilter == "science":
		for index, post := range fullPosts {
			if contains(post.CategoryArr, appliedFilter) {
				filteredPosts = append(filteredPosts, fullPosts[index])
			}
		}
	case appliedFilter == "education":
		for index, post := range fullPosts {
			if contains(post.CategoryArr, appliedFilter) {
				filteredPosts = append(filteredPosts, fullPosts[index])
			}
		}
	case appliedFilter == "sports":
		for index, post := range fullPosts {
			if contains(post.CategoryArr, appliedFilter) {
				filteredPosts = append(filteredPosts, fullPosts[index])
			}
		}
	case appliedFilter == "lifehacks":
		for index, post := range fullPosts {
			if contains(post.CategoryArr, appliedFilter) {
				filteredPosts = append(filteredPosts, fullPosts[index])
			}
		}
	case appliedFilter == "mine":
		// Check if the user is logged-in. If cookie is empty, redirect to the index page.
		cookie, err := r.Cookie("session")
		receivedUUID := cookie.Value
		matchedUsername := getUsernameFromUUID(w, receivedUUID)
		if err != nil || matchedUsername == "" {
			fmt.Println("ERROR: Log-in needed to filter posts")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		for index, post := range fullPosts {
			if post.CreatorUsrName == matchedUsername {
				filteredPosts = append(filteredPosts, fullPosts[index])
			}
		}
	default:
		filteredPosts = fullPosts
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
