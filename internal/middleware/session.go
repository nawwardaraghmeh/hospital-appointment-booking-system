package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// delete the session cookie by setting an expired cookie
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "abs_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
}

// parse the session cookie following the format: role:userID:hospitalID
func GetSessionCookie(r *http.Request) (role string, userID int, hospitalID int, ok bool) {
	c, err := r.Cookie("abs_session")
	if err != nil {
		return "", 0, 0, false
	}
	parts := strings.Split(c.Value, ":")
	if len(parts) != 3 {
		return "", 0, 0, false
	}
	uid, err1 := strconv.Atoi(parts[1])
	hid, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil {
		return "", 0, 0, false
	}
	return parts[0], uid, hid, true
}

// protect a route. check if session cookie exists, and has the allowed role
func RequireSession(next http.HandlerFunc, allowedRole string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, _, _, ok := GetSessionCookie(r)
		if !ok || role != allowedRole {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
