package src

import (
	"database/sql"
	"log"
)

func InsertComment(db *sql.DB, comment Comment) {
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
