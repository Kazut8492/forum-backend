package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

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
