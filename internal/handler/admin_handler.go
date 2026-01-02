package handler

import (
	"abs/internal/repository"
	"html/template"
	"net/http"
)

// Admin login page
func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	db, _ := repository.OpenDB()
	defer db.Close()

	// Cities
	rows, _ := db.Query("SELECT id, name FROM city")
	defer rows.Close()
	type City struct {
		ID   int
		Name string
	}
	var cities []City
	for rows.Next() {
		var c City
		rows.Scan(&c.ID, &c.Name)
		cities = append(cities, c)
	}

	// Hospitals
	hrows, _ := db.Query("SELECT id, name, city_id FROM hospital")
	defer hrows.Close()
	type Hospital struct {
		ID, CityID int
		Name       string
	}
	var hospitals []Hospital
	for hrows.Next() {
		var h Hospital
		hrows.Scan(&h.ID, &h.Name, &h.CityID)
		hospitals = append(hospitals, h)
	}

	tmpl := template.Must(template.ParseFiles("templates/admin_login.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Cities":    cities,
		"Hospitals": hospitals,
	})
}

// Admin dashboard
func AdminPage(w http.ResponseWriter, r *http.Request) {
	db, _ := repository.OpenDB()
	defer db.Close()

	type Timeslot struct {
		ID         int
		Department string
		Doctor     string
		Room       string
		StartTime  string
		Duration   int
		IsBooked   bool
		Patient    string
	}
	rows, _ := db.Query(`
		SELECT t.id, d.name, t.doctor, t.room, t.start_time, t.duration, t.is_booked,
		       IFNULL(a.patient_name,'')
		FROM timeslot t
		JOIN department d ON t.department_id=d.id
		LEFT JOIN appointment a ON t.id=a.timeslot_id
	`)
	defer rows.Close()
	var timeslots []Timeslot
	for rows.Next() {
		var t Timeslot
		var booked int
		rows.Scan(&t.ID, &t.Department, &t.Doctor, &t.Room, &t.StartTime, &t.Duration, &booked, &t.Patient)
		t.IsBooked = booked == 1
		timeslots = append(timeslots, t)
	}

	// Departments for add slot dropdown
	drows, _ := db.Query("SELECT id, name FROM department")
	defer drows.Close()
	type Department struct {
		ID   int
		Name string
	}
	var depts []Department
	for drows.Next() {
		var d Department
		drows.Scan(&d.ID, &d.Name)
		depts = append(depts, d)
	}

	tmpl := template.Must(template.ParseFiles("templates/admin.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Timeslots":   timeslots,
		"Departments": depts,
	})
}

// Add slot
func AddSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		AdminPage(w, r)
		return
	}
	db, _ := repository.OpenDB()
	defer db.Close()

	deptID := r.FormValue("department_id")
	doctor := r.FormValue("doctor")
	room := r.FormValue("room")
	start := r.FormValue("start_time")
	dur := r.FormValue("duration")

	db.Exec(`INSERT INTO timeslot(department_id, doctor, room, start_time, duration, is_booked)
		VALUES(?,?,?,?,?,0)`, deptID, doctor, room, start, dur)

	// Show dashboard instantly
	AdminPage(w, r)
}
