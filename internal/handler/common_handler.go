package handler

import (
	"abs/internal/middleware"
	"html/template"
	"net/http"
)

// HomePage serves the main landing page for both patients and admins
func HomePage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, nil)
}

// Logout destroys the session cookie and redirects back home
func Logout(w http.ResponseWriter, r *http.Request) {
	middleware.ClearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
