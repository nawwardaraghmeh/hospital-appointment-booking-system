package middleware

import (
	"net/http"
	"strings"
)

// ClearSessionCookie instructs the browser to immediately expire the session cookie
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

// GetSessionCookie parses the session cookie into its three parts: role, userID, hospitalID
func GetSessionCookie(r *http.Request) (role, id, hospID string, ok bool) {
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

// RequireSession is a middleware that protects routes, allowing only the specified role through
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
