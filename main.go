package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type TestComment struct {
	ID          int
	Title       string
	Description string
}

type TestPost struct {
	ID          int
	Title       string
	Description string
	Comments    []*TestComment
}

var testPosts []TestPost

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	// query := url.Values{
	// 	"id": []string{},
	// }
	// fmt.Println(query.Encode())

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	sql_table := `
	CREATE TABLE IF NOT EXISTS post(
		Id TEXT NOT NULL PRIMARY KEY,
		Name TEXT,
		Phone TEXT,
		InsertedDatetime DATETIME
	);
	`
	_, err = db.Exec(sql_table)
	if err != nil {
		log.Fatal(err)
	}

	testPosts = append(testPosts, TestPost{ID: 1, Title: "Title", Description: "Description"})
	testPosts = append(testPosts, TestPost{ID: 2, Title: "Title", Description: "Description"})
	testPosts = append(testPosts, TestPost{ID: 3, Title: "Title", Description: "Description"})
	testPosts = append(testPosts, TestPost{ID: 4, Title: "Title", Description: "Description"})
	testPosts = append(testPosts, TestPost{ID: 5, Title: "Title", Description: "Description"})
	testPosts = append(testPosts, TestPost{ID: 6, Title: "Title", Description: "Description"})
	testPosts = append(testPosts, TestPost{ID: 7, Title: "Title", Description: "Description"})
	testPosts[0].Comments = append(testPosts[0].Comments, &TestComment{ID: 1, Title: "CommentTitle", Description: "CommentDescription"})

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/write", writeHandler)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
	http.ListenAndServe(":8888", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}
	if err := tpl.ExecuteTemplate(w, "index.html", testPosts); err != nil {
		w.WriteHeader(500)
		return
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	postIDstr := r.FormValue("id")
	postID, err := strconv.Atoi(postIDstr)
	// fmt.Println(postID)
	if err != nil {
		log.Fatal(err)
	}
	if r.URL.Path != "/post" {
		w.WriteHeader(404)
		return
	}
	if 1 <= postID && postID <= len(testPosts) {
		if err := tpl.ExecuteTemplate(w, "post.html", testPosts[postID-1]); err != nil {
			w.WriteHeader(500)
			return
		}
	} else {
		w.WriteHeader(400)
		return
	}
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var newComment TestComment
	postIDstr := r.FormValue("id")
	postID, err := strconv.Atoi(postIDstr)
	// fmt.Println(postID)
	if err != nil {
		log.Fatal(err)
	}
	newComment.ID = len(testPosts[postID-1].Comments) + 1
	newComment.Title = r.Form["commentTitle"][0]
	newComment.Description = r.Form["commentDescription"][0]
	// fmt.Println(newComment)
	testPosts[postID-1].Comments = append(testPosts[postID-1].Comments, &newComment)
	http.Redirect(w, r, "/post?id="+postIDstr, http.StatusFound)
}
