package handler

import (
	"abs/internal/repository"
	"abs/internal/service"
	"html/template"
	"net/http"
	"strconv"
)

// in a real system, this would be in an environment variable
const AdminSecret = "HOSPITAL2026"

func RegisterAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")
	fullName := r.FormValue("full_name")
	role := r.FormValue("role")

	if role == "admin" {
		secret := r.FormValue("admin_code")
		if secret != AdminSecret {
			http.Error(w, "Invalid Admin Registration Code", http.StatusUnauthorized)
			return
		}
	}

	// 1. hash the password
	hash, err := service.HashPassword(password)
	if err != nil {
		http.Error(w, "Error processing password", 500)
		return
	}

	// 2. open DB and Save
	db, _ := repository.OpenDB()
	defer db.Close()

	var hospitalID interface{}
	if role == "admin" {
		hospitalID, _ = strconv.Atoi(r.FormValue("hospital_id"))
	}

	err = repository.CreateUser(db, username, hash, role, fullName, hospitalID)
	if err != nil {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// RegisterPage renders the signup form
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	db, _ := repository.OpenDB()
	defer db.Close()

	rows, _ := db.Query("SELECT id, name FROM hospital")
	var hospitals []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		hospitals = append(hospitals, map[string]interface{}{
			"ID":   id,
			"Name": name,
		})
	}

	tmpl, _ := template.ParseFiles("templates/register.html")
	tmpl.Execute(w, map[string]interface{}{
		"Hospitals": hospitals,
	})
}

// LoginAction checks if user exists then redirects
func LoginAction(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	db, _ := repository.OpenDB()
	defer db.Close()

	// 1. find the user in our new 'users' table
	var id int
	var hash, role string
	var hospID int
	err := db.QueryRow("SELECT id, password_hash, role, IFNULL(hospital_id, 0) FROM users WHERE username = ?",
		username).Scan(&id, &hash, &role, &hospID)

	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// 2. use the Service to verify the hashed password
	if !service.CheckPasswordHash(password, hash) {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// 3. set the session. it stores role, userID, and hospitalID in the cookie
	cookieValue := role + ":" + strconv.Itoa(id) + ":" + strconv.Itoa(hospID)
	http.SetCookie(w, &http.Cookie{
		Name:     "abs_session",
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
	})

	// 4. redirect based on role
	if role == "admin" {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		deptID := r.FormValue("department_id")
		if deptID != "" {
			http.Redirect(w, r, "/patient/slots?department_id="+deptID, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/patient/slots", http.StatusSeeOther)
		}
	}
}
