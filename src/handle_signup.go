package src

import (
	"html/template"
	"net/http"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	if r.URL.Path != "/signup" {
		w.WriteHeader(404)
		return
	}

	if err := tpl.ExecuteTemplate(w, "signup.html", nil); err != nil {
		w.WriteHeader(500)
		return
	}
}
