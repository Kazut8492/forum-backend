package data

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func PasswordEncrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func CreateTables(db *sql.DB) {
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

func InsertPosts(db *sql.DB, posts []Post) {
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
		err = rows.Scan(&post.Id, &post.Title, &post.Content)
		if err != nil {
			panic(err.Error())
		}
		result = append(result, post)
	}
	return result
}

func InsertComments(db *sql.DB, comment Comment) {
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
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.Title, &comment.Content)
		if err != nil {
			panic(err)
		}
		result = append(result, comment)
	}
	return result
}

func InsertUser(db *sql.DB, newUser User) {
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

func ReadUser(db *sql.DB, loginUser User) User {
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

func InitiateSession(w http.ResponseWriter, r *http.Request, db *sql.DB, user User) {
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
