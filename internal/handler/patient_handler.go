package handler

import (
	"abs/internal/repository"
	"html/template"
	"net/http"
)

// --- Structs used throughout ---
type City struct {
	ID   int
	Name string
}

type Hospital struct {
	ID     int
	CityID int
	Name   string
}

type Department struct {
	ID         int
	HospitalID int
	Name       string
}

type Slot struct {
	ID        int
	Doctor    string
	Room      string
	StartTime string
	Duration  int
}

// --- Index Page ---
func IndexPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, nil)
}

// --- Logout Page ---
func LogoutPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// --- Patient Login Page ---
func PatientLoginPage(w http.ResponseWriter, r *http.Request) {
	db, _ := repository.OpenDB()
	defer db.Close()

	// Load Cities
	cities := []City{}
	rows, _ := db.Query("SELECT id, name FROM city")
	defer rows.Close()
	for rows.Next() {
		var c City
		rows.Scan(&c.ID, &c.Name)
		cities = append(cities, c)
	}

	// Load Hospitals
	hospitals := []Hospital{}
	hrows, _ := db.Query("SELECT id, name, city_id FROM hospital")
	defer hrows.Close()
	for hrows.Next() {
		var h Hospital
		hrows.Scan(&h.ID, &h.Name, &h.CityID)
		hospitals = append(hospitals, h)
	}

	// Load Departments dynamically
	departments := []Department{}
	drows, _ := db.Query("SELECT id, name, hospital_id FROM department")
	defer drows.Close()
	for drows.Next() {
		var d Department
		drows.Scan(&d.ID, &d.Name, &d.HospitalID)
		departments = append(departments, d)
	}

	// Render template
	t := template.Must(template.ParseFiles("templates/patient_login.html"))
	t.Execute(w, map[string]interface{}{
		"Cities":      cities,
		"Hospitals":   hospitals,
		"Departments": departments, // <-- now always fetched from DB
	})
}

// --- Patient Slots Page ---
func PatientSlotsPage(w http.ResponseWriter, r *http.Request) {
	deptID := r.URL.Query().Get("department_id")
	if deptID == "" {
		http.Error(w, "No department selected", 400)
		return
	}

	db, _ := repository.OpenDB()
	defer db.Close()

	// Fetch slots for the selected department only
	rows, _ := db.Query(`
		SELECT id, doctor, room, start_time, duration
		FROM timeslot
		WHERE department_id=? AND is_booked=0
	`, deptID)
	defer rows.Close()

	slots := []Slot{}
	for rows.Next() {
		var s Slot
		rows.Scan(&s.ID, &s.Doctor, &s.Room, &s.StartTime, &s.Duration)
		slots = append(slots, s)
	}

	t := template.Must(template.ParseFiles("templates/patient_slots.html"))
	t.Execute(w, map[string]interface{}{"Slots": slots, "DepartmentID": deptID})
}

// --- Book Appointment ---
func BookAppointment(w http.ResponseWriter, r *http.Request) {
	slotID := r.URL.Query().Get("slot_id")
	success := false

	db, _ := repository.OpenDB()
	defer db.Close()

	if r.Method == http.MethodPost {
		db.Exec(`INSERT INTO appointment(timeslot_id, patient_name, patient_id, age, phone, email, symptoms)
			VALUES(?,?,?,?,?,?,?)`,
			slotID,
			r.FormValue("name"),
			r.FormValue("pid"),
			r.FormValue("age"),
			r.FormValue("phone"),
			r.FormValue("email"),
			r.FormValue("symptoms"),
		)
		db.Exec(`UPDATE timeslot SET is_booked=1 WHERE id=?`, slotID)
		success = true
	}

	t := template.Must(template.ParseFiles("templates/book.html"))
	t.Execute(w, map[string]interface{}{
		"SlotID":  slotID,
		"Success": success,
	})
}
