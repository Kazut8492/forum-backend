package src

import (
	"html/template"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	if r.URL.Path != "/login" {
		w.WriteHeader(404)
		return
	}

	if err := tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		w.WriteHeader(500)
		return
	}
}
