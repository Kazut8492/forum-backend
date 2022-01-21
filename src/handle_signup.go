package src

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	if r.URL.Path != "/signup" {
		w.WriteHeader(404)
		return
	}

	db, err := sql.Open("sqlite3", "./example.db")
	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err.Error())
		log.Fatal(1)
	}
	defer db.Close()

	//Get warnings
	warningRows, err := db.Query("SELECT warning_type FROM warning")
	if err != nil {
		w.WriteHeader(500)
		log.Fatal(err.Error())
	}
	defer warningRows.Close()
	var warnings []string
	for warningRows.Next() {
		var warning string
		err = warningRows.Scan(&warning)
		if err != nil {
			w.WriteHeader(500)
			panic(err.Error())
		}
		warnings = append(warnings, warning)
	}

	//Reset Warnings
	db.Exec("DELETE FROM warning")

	if err := tpl.ExecuteTemplate(w, "signup.html", warnings); err != nil {
		w.WriteHeader(500)
		return
	}
}
