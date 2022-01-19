package src

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	cookie, err := r.Cookie("session")
	if err != nil {
		fmt.Println("ERROR: Failed to get cookie info")
		http.Redirect(w, r, "/", http.StatusFound)
		return
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
