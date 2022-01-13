package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateTables(db *sql.DB) {
	dbTables := []string{
		`CREATE TABLE IF NOT EXISTS post (
			"post_id"			INTEGER NOT NULL UNIQUE,
			"title"				TEXT NOT NULL,
			"content"			TEXT NOT NULL,
			"category_str"		TEXT NOT NULL,
			"creator_username"	TEXT NOT NULL,
			PRIMARY KEY("post_id" AUTOINCREMENT)
			FOREIGN KEY("creator_username") REFERENCES "USER"("username")
		)`,

		`CREATE TABLE IF NOT EXISTS comment (
			"comment_id"		INTEGER NOT NULL UNIQUE,
			"post_id"			INTEGER NOT NULL,
			"title"				TEXT NOT NULL,
			"content"			TEXT NOT NULL,
			"creator_username"	TEXT NOT NULL,
			PRIMARY KEY("comment_id" AUTOINCREMENT),
			FOREIGN KEY("post_id") REFERENCES "POST"("post_id")
			FOREIGN KEY("creator_username") REFERENCES "USER"("username")
		)`,

		`CREATE TABLE IF NOT EXISTS user (
			"user_id"		INTEGER NOT NULL UNIQUE,
			"username"		TEXT NOT NULL UNIQUE,
			"user_email"	TEXT NOT NULL UNIQUE,
			"user_pass"		TEXT NOT NULL,
			PRIMARY KEY("user_id" AUTOINCREMENT)
		)`,

		`CREATE TABLE IF NOT EXISTS session (
			"session_id"	INTEGER NOT NULL UNIQUE,
			"datetime"		DATETIME DEFAULT CURRENT_TIMESTAMP,
			"username"		TEXT NOT NULL UNIQUE,
			"uuid"			TEXT NOT NULL,
			PRIMARY KEY("session_id" AUTOINCREMENT),
			FOREIGN KEY("username") REFERENCES "USER"("username")
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

func InsertPost(db *sql.DB, post Post) {
	// comment_id shall be inserted automatically, also be careful to match VALUES

	categoryStr := strings.Join(post.CategoryArr, ",")

	statement, err := db.Prepare(`
		INSERT INTO post (
			title,
			content,
			category_str,
			creator_username
		) VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer statement.Close()
	// number of variables have to be matched with INSERTed variables
	statement.Exec(post.Title, post.Content, categoryStr, post.CreatorUsrName)
}

func ReadPosts(db *sql.DB) []Post {
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
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CategoryStr, &post.CreatorUsrName)
		if err != nil {
			panic(err.Error())
		}
		result = append(result, post)
	}
	for index, post := range result {
		result[index].CategoryArr = strings.Split(post.CategoryStr, ",")
	}

	return result
}

func InsertComments(db *sql.DB, comment Comment) {
	// comment_id shall be inserted automatically, also be careful to match VALUES
	statement, err := db.Prepare(`
		INSERT INTO comment (
			post_id,
			title,
			content,
			creator_username
		) VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer statement.Close()
	// number of variables have to be matched with INSERTed variables
	statement.Exec(comment.PostId, comment.Title, comment.Content, comment.CreatorUsrName)
}

func ReadComments(db *sql.DB, postId int) []Comment {
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
		err = rows.Scan(&comment.ID, &comment.PostId, &comment.Title, &comment.Content, &comment.CreatorUsrName)
		if err != nil {
			panic(err)
		}
		result = append(result, comment)
	}
	return result
}

func InsertUser(db *sql.DB, user User) {
	statement, err := db.Prepare(`
		INSERT INTO user (
			username,
			user_email,
			user_pass
		) VALUES (?, ?, ?)
	`)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	defer statement.Close()
	statement.Exec(user.Username, user.Email, user.Pass)
}

// Commented out becasue there is one-liner way to do it. Pls refer to Login/Signup Submit Handlers
// func ReadUser(w http.ResponseWriter, db *sql.DB, loginUser User) User {
// 	statement, _ := db.Query("SELECT * FROM user WHERE username = ?", loginUser.Username)
// 	// if err != nil {
// 	// 	w.WriteHeader(500)
// 	// 	fmt.Println(err.Error())
// 	// 	log.Fatal(1)
// 	// }
// 	defer statement.Close()

// 	user := User{}
// 	for statement.Next() {
// 		err := statement.Scan(&user.ID, &user.Username, &user.Email, &user.Pass)
// 		if err != nil {
// 			fmt.Println(err.Error())
// 			return User{}
// 		}
// 	}
// 	return user
// }

func getUsernameFromUUID(w http.ResponseWriter, receivedUUID string) string {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	defer db.Close()

	var matchedUsername string
	db.QueryRow("SELECT username FROM session WHERE uuid = ?", receivedUUID).Scan(&matchedUsername)
	return matchedUsername
}

func InitiateSession(w http.ResponseWriter, r *http.Request, db *sql.DB, user User) {
	uuid := uuid.New()
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	db.Exec("DELETE FROM session WHERE user_id = ?", user.ID)

	cookie := http.Cookie{
		Name:    "session",
		Value:   uuid.String(),
		Expires: expiration,
		Secure:  true,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)

	statement, err := db.Prepare("INSERT INTO session (username ,uuid) VALUES (?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("ERROR: Failed to insert session")
		log.Fatal(1)
	}
	defer statement.Close()
	statement.Exec(user.Username, uuid)
	http.SetCookie(w, &cookie)
}
