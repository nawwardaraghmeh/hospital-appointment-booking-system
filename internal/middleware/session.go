package middleware

import (
	"net/http"
	"strings"
)

func SetSessionCookie(w http.ResponseWriter, userType string, id string) {
	cookie := http.Cookie{
		Name:     "abs_session",
		Value:    userType + ":" + id,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
}

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

func GetSessionCookie(r *http.Request) (userType string, id string, ok bool) {
	c, err := r.Cookie("abs_session")
	if err != nil {
		return "", "", false
	}

	parts := strings.SplitN(c.Value, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

func RequireSession(next http.HandlerFunc, allowedUserType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userType, _, ok := GetSessionCookie(r)
		if !ok || userType != allowedUserType {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
