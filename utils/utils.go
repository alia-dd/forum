package utils

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"gitea.kood.tech/jyrkikarhunen/forum/models"
)

var Tpl *template.Template

func InitializeTemplate() {
	var TplErr error
	Tpl, TplErr = template.New("").ParseGlob("templates/partials/*.html")
	if TplErr != nil {
		log.Fatal(TplErr)
	}
	Tpl, TplErr = Tpl.ParseFiles("templates/base.html")
	if TplErr != nil {
		log.Fatal(TplErr)
	}

}

func RenderTemplate(w http.ResponseWriter, statusCode int, templateString string, payload any) {
	t, tmpErr := Tpl.Clone()
	if tmpErr != nil {
		http.Error(w, tmpErr.Error(), http.StatusInternalServerError)
		return
	}
	if _, tmpErr := t.ParseFiles("templates/" + templateString + ".html"); tmpErr != nil {
		log.Println("parse template error :", tmpErr)
		serverErrorHandler(w)
		return
	}
	renderExecuted(w, t, statusCode, payload)
}

func serverErrorHandler(w http.ResponseWriter) {
	t, tmpErr := Tpl.Clone()
	if tmpErr != nil {
		log.Println("clone err::", tmpErr)
		http.Error(w, tmpErr.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := t.ParseFiles("templates/" + "error.html"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderExecuted(w, t, http.StatusOK, &models.ErrorStruct{
		Error:   "500",
		ErrorMs: "Internal Server Error",
	})

}

func renderExecuted(w http.ResponseWriter, t *template.Template, statusCode int, payload any) {
	var buf bytes.Buffer

	if err := t.ExecuteTemplate(&buf, "base", payload); err != nil {
		// handleError()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	buf.WriteTo(w)
}

func RenderPartial(w http.ResponseWriter, statusCode int, templateName string, payload any) {
	t, tmpErr := Tpl.Clone()
	if tmpErr != nil {
		http.Error(w, tmpErr.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, templateName, payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	buf.WriteTo(w)
}

func RedirectTologin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/user/login")
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
