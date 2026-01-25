package middleware

import (
	"net/http"
	"strings"
)

// ClearSessionCookie logs the user out by instructing the browser to delete the cookie
func ClearSessionCookie(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "abs_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

// GetSessionCookie retrieves and parses the session data
func GetSessionCookie(r *http.Request) (role string, id string, hospID string, ok bool) {
	c, err := r.Cookie("abs_session")
	if err != nil {
		return "", "", "", false
	}

	parts := strings.Split(c.Value, ":")
	if len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

// RequireSession wraps a route handler and ensures only users with the correct role can enter
func RequireSession(next http.HandlerFunc, allowedUserType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, _, _, ok := GetSessionCookie(r)
		if !ok || role != allowedUserType {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
