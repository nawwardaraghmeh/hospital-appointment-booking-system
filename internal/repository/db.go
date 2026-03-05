package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// OpenDB opens an SQLite connection to the given path
func OpenDB(dbPath string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("could not create db directory: %w", err)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("could not open database at %s: %w", dbPath, err)
	}
	return db, nil
}

// InitDB creates all tables if they do not already exist
func InitDB(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS city(
			id   INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS hospital(
			id      INTEGER PRIMARY KEY AUTOINCREMENT,
			name    TEXT,
			city_id INTEGER,
			UNIQUE(name, city_id),
			FOREIGN KEY(city_id) REFERENCES city(id)
		)`,
		`CREATE TABLE IF NOT EXISTS department(
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT,
			hospital_id INTEGER,
			UNIQUE(name, hospital_id),
			FOREIGN KEY(hospital_id) REFERENCES hospital(id)
		)`,
		`CREATE TABLE IF NOT EXISTS timeslot(
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			department_id INTEGER,
			doctor        TEXT,
			room          TEXT,
			start_time    TEXT,
			duration      INTEGER,
			is_booked     INTEGER DEFAULT 0,
			FOREIGN KEY(department_id) REFERENCES department(id)
		)`,
		`CREATE TABLE IF NOT EXISTS appointment(
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			timeslot_id  INTEGER,
			patient_name TEXT,
			patient_id   TEXT,
			age          INTEGER,
			phone        TEXT,
			email        TEXT,
			symptoms     TEXT,
			FOREIGN KEY(timeslot_id) REFERENCES timeslot(id)
		)`,
		`CREATE TABLE IF NOT EXISTS users(
			id            INTEGER PRIMARY KEY AUTOINCREMENT,
			username      TEXT UNIQUE,
			password_hash TEXT,
			role          TEXT,
			full_name     TEXT,
			hospital_id   INTEGER,
			FOREIGN KEY(hospital_id) REFERENCES hospital(id)
		)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("schema error: %w", err)
		}
	}
	return nil
}
