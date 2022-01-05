package src

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
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

	certainPost := Post{posts[postID-1].Id, posts[postID-1].Title, posts[postID-1].Content, comments}

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

func WriteHandler(w http.ResponseWriter, r *http.Request) {
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
