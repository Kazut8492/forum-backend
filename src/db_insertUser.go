package src

import (
	"database/sql"
	"fmt"
	"log"
)

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
