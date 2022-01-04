package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Session struct {
	userId     int
	username   string
	sessionId  int
	expiration time.Time
}

type User struct {
	User_ID  int
	Username string
	Pass     string
}

type Comment struct {
	Id      int
	PostId  int
	Title   string
	Content string
}

type Post struct {
	Id       int
	Title    string
	Content  string
	Comments []Comment
}

// var testPosts []Post

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
		log.Fatal(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	createTables(db)
	testPosts := []Post{
		{Title: "Title1", Content: "Content1"},
		{Title: "Title2", Content: "Content2"},
		{Title: "Title3", Content: "Content3"},
		{Title: "Title4", Content: "Content4"},
		{Title: "Title5", Content: "Content5"},
		{Title: "Title6", Content: "Content6"},
		{Title: "Title7", Content: "Content7"},
	}
	insertPosts(db, testPosts)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/signup-submit", signupSubmitHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/login-submit", loginSubmitHandler)
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates/"))))
	http.ListenAndServe(":8888", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	posts := readPosts(db)

	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}
	if err := tpl.ExecuteTemplate(w, "index.html", posts); err != nil {
		w.WriteHeader(500)
		return
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	posts := readPosts(db)

	postIDstr := r.FormValue("id")
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		log.Fatal(err.Error())
	}

	comments := readComments(db, postID)

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

func writeHandler(w http.ResponseWriter, r *http.Request) {
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
	insertComments(db, newComment)
	// fmt.Println(newComment)
	// testPosts[postID-1].Comments = append(testPosts[postID-1].Comments, &newComment)
	http.Redirect(w, r, "/post?id="+postIDstr, http.StatusFound)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		w.WriteHeader(404)
		return
	}
	if err := tpl.ExecuteTemplate(w, "signup.html", nil); err != nil {
		w.WriteHeader(500)
		return
	}
}

func signupSubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	userPass := r.FormValue("userPass")
	encryptedUserPass, err := PasswordEncrypt(userPass)
	if err != nil {
		fmt.Println("ERROR: Failed to encrypt passworrd")
		log.Fatal(1)
	}
	newUser := User{}
	newUser.Username = username
	newUser.Pass = encryptedUserPass
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	insertUser(db, newUser)
	http.Redirect(w, r, "/", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		w.WriteHeader(404)
		return
	}
	if err := tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		w.WriteHeader(500)
		return
	}
}

func loginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	userPass := r.FormValue("userPass")
	loginUser := User{}
	loginUser.Username = username
	loginUser.Pass = userPass
	// fmt.Printf("Input login user info: %v", loginUser)
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	userData := readUser(db, loginUser)
	if (userData == User{}) {
		fmt.Println("log in failed, empty data returned")
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		if err := CompareHashAndPassword(userData.Pass, userPass); err != nil {
			fmt.Println("log in failed, password not matched")
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			fmt.Println("log in successed")
			initiateSession(w, r, db, userData)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func createTables(db *sql.DB) {
	dbTables := []string{
		`CREATE TABLE IF NOT EXISTS post (
			"post_id"	INTEGER NOT NULL UNIQUE,
			"title"		TEXT NOT NULL,
			"content"	TEXT NOT NULL,
			PRIMARY KEY("post_id" AUTOINCREMENT)
		)`,

		`CREATE TABLE IF NOT EXISTS comment (
			"comment_id"	INTEGER NOT NULL UNIQUE,
			"post_id"		INTEGER NOT NULL,
			"title"			TEXT NOT NULL,
			"content"		TEXT NOT NULL,
			PRIMARY KEY("comment_id" AUTOINCREMENT),
			FOREIGN KEY("post_id") REFERENCES "POST"("post_id")
		)`,

		`CREATE TABLE IF NOT EXISTS user (
			"user_id"		INTEGER UNIQUE NOT NULL,
			"username"		TEXT NOT NULL UNIQUE,
			"user_pass"		TEXT NOT NULL,
			PRIMARY KEY("user_id" AUTOINCREMENT)
		)`,

		`CREATE TABLE IF NOT EXISTS session (
			"session_id"	INTEGER NOT NULL UNIQUE,
			"datetime"		DATETIME DEFAULT CURRENT_TIMESTAMP,
			"user_id"		INTEGER NOT NULL,
			"username"		TEXT NOT NULL UNIQUE,
			"uuid"			TEXT NOT NULL,
			PRIMARY KEY("session_id" AUTOINCREMENT),
			FOREIGN KEY("user_id") REFERENCES "USER"("user_id")
		)`,
	}
	for _, table := range dbTables {
		statement, err := db.Prepare(table)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer statement.Close()
		statement.Exec()
	}
}

func insertPosts(db *sql.DB, posts []Post) {
	// comment_id shall be inserted automatically, also be careful to match VALUES
	db_storePosts := `
		INSERT INTO post (
			title,
			content
		) VALUES (?, ?)
	`
	statement, err := db.Prepare(db_storePosts)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer statement.Close()
	for _, post := range posts {
		// number of variables have to be matched with INSERTed variables
		statement.Exec(post.Title, post.Content)
	}
}

func readPosts(db *sql.DB) []Post {
	db_readPosts := `
		SELECT * FROM post
		ORDER BY post_id
	`

	rows, err := db.Query(db_readPosts)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var result []Post
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content)
		if err != nil {
			panic(err.Error())
		}
		result = append(result, post)
	}
	return result
}

func insertComments(db *sql.DB, comment Comment) {
	// comment_id shall be inserted automatically, also be careful to match VALUES
	db_storeComments := `
		INSERT INTO comment (
			post_id,
			title,
			content
		) VALUES (?, ?, ?)
	`
	statement, err := db.Prepare(db_storeComments)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer statement.Close()
	// number of variables have to be matched with INSERTed variables
	statement.Exec(comment.PostId, comment.Title, comment.Content)
}

func readComments(db *sql.DB, postId int) []Comment {
	db_readComments := `
		SELECT * FROM comment WHERE post_id = ?
		ORDER BY comment_id
	`

	rows, err := db.Query(db_readComments, postId)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var result []Comment
	for rows.Next() {
		comment := Comment{}
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.Title, &comment.Content)
		if err != nil {
			panic(err)
		}
		result = append(result, comment)
	}
	return result
}

func insertUser(db *sql.DB, newUser User) {
	db_storeUser := `
		INSERT INTO user (
			username,
			user_pass
		) VALUES (?, ?)
	`
	statement, err := db.Prepare(db_storeUser)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer statement.Close()
	statement.Exec(newUser.Username, newUser.Pass)
}

func readUser(db *sql.DB, loginUser User) User {
	db_readUser := `
		SELECT * FROM user WHERE username = ?
	`

	row, err := db.Query(db_readUser, loginUser.Username)
	if err != nil {
		fmt.Println("Username not found in the database")
		return User{}
	}
	defer row.Close()

	// var result []Comment
	user := User{}
	for row.Next() {
		err = row.Scan(&user.User_ID, &user.Username, &user.Pass)
		if err != nil {
			fmt.Println(err.Error())
			return User{}
		}
	}
	return user
}

func initiateSession(w http.ResponseWriter, r *http.Request, db *sql.DB, user User) {
	uuid := uuid.New()
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	db.Exec("DELETE FROM session WHERE user_id = ?", user.User_ID)

	cookie := http.Cookie{
		Name:    "session",
		Value:   uuid.String(),
		Expires: expiration,
		Secure:  true,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	statement, err := db.Prepare("INSERT INTO session (user_id, username ,uuid, datetime) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("ERROR: Failed to insert session")
		log.Fatal(1)
	}
	defer statement.Close()
	statement.Exec(user.User_ID, user.Username, uuid, expiration)
	http.SetCookie(w, &cookie)
}
