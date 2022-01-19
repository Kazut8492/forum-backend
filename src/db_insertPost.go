package src

import (
	"database/sql"
	"log"
	"strings"
)

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
