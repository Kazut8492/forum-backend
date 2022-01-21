package src

import (
	"database/sql"
	"log"
)

func CreateTables(db *sql.DB) {
	dbTables := []string{
		`CREATE TABLE IF NOT EXISTS post (
			"post_id"			INTEGER NOT NULL UNIQUE,
			"title"				TEXT NOT NULL,
			"content"			TEXT NOT NULL,
			"category_str"		TEXT NOT NULL,
			"creator_username"	TEXT NOT NULL,
			PRIMARY KEY("post_id" AUTOINCREMENT),
			FOREIGN KEY("creator_username") REFERENCES "USER"("username")
		)`,

		`CREATE TABLE IF NOT EXISTS comment (
			"comment_id"		INTEGER NOT NULL UNIQUE,
			"post_id"			INTEGER NOT NULL,
			"title"				TEXT NOT NULL,
			"content"			TEXT NOT NULL,
			"creator_username"	TEXT NOT NULL,
			PRIMARY KEY("comment_id" AUTOINCREMENT),
			FOREIGN KEY("post_id") REFERENCES "POST"("post_id"),
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

		`CREATE TABLE IF NOT EXISTS like (
			"like_id"			INTEGER NOT NULL UNIQUE,
			"post_id"			INTEGER NOT NULL,
			"comment_id"		INTEGER DEFAULT '0',
			"creator_username"	TEXT NOT NULL,
			PRIMARY KEY("like_id" AUTOINCREMENT),
			FOREIGN KEY("creator_username") REFERENCES "USER"("username"),
			FOREIGN KEY("post_id") REFERENCES "POST"("post_id"),
			FOREIGN KEY("comment_id") REFERENCES "COMMENT"("comment_id")
		)`,

		`CREATE TABLE IF NOT EXISTS dislike (
			"dislike_id"		INTEGER NOT NULL UNIQUE,
			"post_id"			INTEGER NOT NULL,
			"comment_id"		INTEGER DEFAULT '0',
			"creator_username"	TEXT NOT NULL,
			PRIMARY KEY("dislike_id" AUTOINCREMENT),
			FOREIGN KEY("creator_username") REFERENCES "USER"("username"),
			FOREIGN KEY("post_id") REFERENCES "POST"("post_id"),
			FOREIGN KEY("comment_id") REFERENCES "COMMENT"("comment_id")
		)`,

		`CREATE TABLE IF NOT EXISTS warning (
			"warning_id"		INTEGER NOT NULL UNIQUE,
			"warning_type"		TEXT NOT NULL UNIQUE,
			PRIMARY KEY("warning_id" AUTOINCREMENT)
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
