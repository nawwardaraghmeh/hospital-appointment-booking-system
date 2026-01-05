package handler

import (
	"abs/internal/entity"
	"abs/internal/middleware"
	"abs/internal/repository"
	"html/template"
	"net/http"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, nil)
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	middleware.ClearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func PatientLoginAction(w http.ResponseWriter, r *http.Request) {
	deptID := r.FormValue("department_id")
	patientID := r.FormValue("pid")

	if deptID == "" || patientID == "" {
		http.Redirect(w, r, "/patient/login", http.StatusSeeOther)
		return
	}

	middleware.SetSessionCookie(w, "patient", patientID)

	target := "/patient/slots?department_id=" + deptID
	http.Redirect(w, r, target, http.StatusSeeOther)
}

func PatientLoginPage(w http.ResponseWriter, r *http.Request) {
	db, _ := repository.OpenDB()
	defer db.Close()

	var cities []entity.City
	rows, _ := db.Query("SELECT id, name FROM city")
	defer rows.Close()
	for rows.Next() {
		var c entity.City
		rows.Scan(&c.ID, &c.Name)
		cities = append(cities, c)
	}

	var hospitals []entity.Hospital
	hrows, _ := db.Query("SELECT id, name, city_id FROM hospital")
	defer hrows.Close()
	for hrows.Next() {
		var h entity.Hospital
		hrows.Scan(&h.ID, &h.Name, &h.CityID)
		hospitals = append(hospitals, h)
	}

	var departments []entity.Department
	drows, _ := db.Query("SELECT id, name, hospital_id FROM department")
	defer drows.Close()
	for drows.Next() {
		var d entity.Department
		drows.Scan(&d.ID, &d.Name, &d.HospitalID)
		departments = append(departments, d)
	}

	t := template.Must(template.ParseFiles("templates/patient_login.html"))
	t.Execute(w, map[string]interface{}{
		"Cities":      cities,
		"Hospitals":   hospitals,
		"Departments": departments,
	})
}

func PatientSlotsPage(w http.ResponseWriter, r *http.Request) {
	deptID := r.URL.Query().Get("department_id")
	if deptID == "" {
		http.Error(w, "No department selected", 400)
		return
	}

	middleware.SetSessionCookie(w, "patient", deptID)

	db, _ := repository.OpenDB()
	defer db.Close()

	rows, _ := db.Query(`
		SELECT id, doctor, room, start_time, duration
		FROM timeslot
		WHERE department_id=? AND is_booked=0
	`, deptID)
	defer rows.Close()

	type Slot struct {
		ID        int
		Doctor    string
		Room      string
		StartTime string
		Duration  int
	}

	var slots []Slot
	for rows.Next() {
		var s Slot
		rows.Scan(&s.ID, &s.Doctor, &s.Room, &s.StartTime, &s.Duration)
		slots = append(slots, s)
	}

	t := template.Must(template.ParseFiles("templates/patient_slots.html"))
	t.Execute(w, map[string]interface{}{
		"Slots":        slots,
		"DepartmentID": deptID,
	})
}

func BookAppointment(w http.ResponseWriter, r *http.Request) {
	slotID := r.URL.Query().Get("slot_id")

	db, _ := repository.OpenDB()
	defer db.Close()

	if r.Method == http.MethodPost {
		var deptID string
		db.QueryRow("SELECT department_id FROM timeslot WHERE id = ?", slotID).Scan(&deptID)

		db.Exec(`
            INSERT INTO appointment (timeslot_id, patient_name, patient_id, age, phone, email, symptoms)
            VALUES (?,?,?,?,?,?,?)`,
			slotID, r.FormValue("name"), r.FormValue("pid"), r.FormValue("age"),
			r.FormValue("phone"), r.FormValue("email"), r.FormValue("symptoms"),
		)
		db.Exec(`UPDATE timeslot SET is_booked=1 WHERE id=?`, slotID)

		http.Redirect(w, r, "/patient/slots?department_id="+deptID, http.StatusSeeOther)
		return
	}

	t := template.Must(template.ParseFiles("templates/book.html"))
	t.Execute(w, map[string]interface{}{
		"SlotID": slotID,
	})
}
