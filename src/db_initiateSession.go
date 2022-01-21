package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

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
		w.WriteHeader(500)
		fmt.Println(err.Error())
		fmt.Println("ERROR: Failed to insert session")
		log.Fatal(1)
	}
	defer statement.Close()
	statement.Exec(user.Username, uuid)
	http.SetCookie(w, &cookie)
}
