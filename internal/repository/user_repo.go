package repository

import (
	"database/sql"
)

// CreateUser saves a new user to the database
func CreateUser(db *sql.DB, username, hash, role, name string, hospitalID interface{}) error {
	_, err := db.Exec(`INSERT INTO users (username, password_hash, role, full_name, hospital_id) 
                       VALUES (?, ?, ?, ?, ?)`,
		username, hash, role, name, hospitalID)
	return err
}

// GetUserByUsername finds a user for authentication
func GetUserByUsername(db *sql.DB, username string) (id int, hash, role string, hospID int, err error) {
	err = db.QueryRow("SELECT id, password_hash, role, IFNULL(hospital_id, 0) FROM users WHERE username = ?",
		username).Scan(&id, &hash, &role, &hospID)
	return
}
