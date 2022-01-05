package src

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	posts := ReadPosts(db)

	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}
	if err := tpl.ExecuteTemplate(w, "index.html", posts); err != nil {
		w.WriteHeader(500)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	if r.URL.Path != "/login" {
		w.WriteHeader(404)
		return
	}
	if err := tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		w.WriteHeader(500)
		return
	}
}

func LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	loginUser := User{}
	loginUser.Username = username
	loginUser.Pass = password
	// fmt.Printf("Input login user info: %v", loginUser)
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	userData := ReadUser(db, loginUser)
	if (userData == User{}) {
		fmt.Println("log in failed, empty data returned")
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		if err := CompareHashAndPassword(userData.Pass, password); err != nil {
			fmt.Println("log in failed, password not matched")
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			fmt.Println("log in successed")
			InitiateSession(w, r, db, userData)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("ERRPR: Failed to get cookie info")
	}
	db.Exec("DELETE FROM session WHERE uuid = ?", cookie.Value)
	cookie = &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusFound)
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
