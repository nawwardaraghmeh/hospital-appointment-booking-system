package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./data/abs.db")
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

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
	db.Exec(`CREATE TABLE IF NOT EXISTS admin(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		emp_id TEXT UNIQUE,
		name TEXT,
		city_id INTEGER,
		hospital_id INTEGER
	)`)
}

func PrintTimeslots() {
	db, _ := OpenDB()
	defer db.Close()

	rows, _ := db.Query(`
        SELECT id, department_id, doctor, room, start_time, duration, is_booked
        FROM timeslot
    `)
	defer rows.Close()

	for rows.Next() {
		var id, deptID, duration, isBooked int
		var doctor, room, start string
		rows.Scan(&id, &deptID, &doctor, &room, &start, &duration, &isBooked)
		fmt.Println(id, deptID, doctor, room, start, duration, isBooked)
	}
}
