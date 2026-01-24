package repository

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// OpenDB opens a connection to the SQLite database
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data/abs.db")
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

// InitDB sets up the database schema
func InitDB() {
	db, _ := OpenDB()
	defer db.Close()

	db.Exec(`CREATE TABLE IF NOT EXISTS city(
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
        name TEXT UNIQUE
    )`)

	db.Exec(`CREATE TABLE IF NOT EXISTS hospital(
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
        name TEXT, 
        city_id INTEGER,
        UNIQUE(name, city_id)
    )`)

	db.Exec(`CREATE TABLE IF NOT EXISTS department(
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
        name TEXT, 
        hospital_id INTEGER,
        UNIQUE(name, hospital_id)
    )`)

	db.Exec(`CREATE TABLE IF NOT EXISTS timeslot(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		department_id INTEGER,
		doctor TEXT,
		room TEXT,
		start_time TEXT,
		duration INTEGER,
		is_booked INTEGER DEFAULT 0
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS appointment(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timeslot_id INTEGER,
		patient_name TEXT,
		patient_id TEXT,
		age INTEGER,
		phone TEXT,
		email TEXT,
		symptoms TEXT
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE,   
    password_hash TEXT,     
    role TEXT,              
    full_name TEXT,
    hospital_id INTEGER,    
    FOREIGN KEY(hospital_id) REFERENCES hospital(id)
)`)
}
