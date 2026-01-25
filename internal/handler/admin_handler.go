package handler

import (
	"abs/internal/entity"
	"abs/internal/middleware"
	"abs/internal/repository"
	"html/template"
	"net/http"
)

// AdminLoginPage loads cities and hospitals for dynamic login dropdowns
func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/admin_login.html"))
	t.Execute(w, nil)
}

// AdminPage displays all appointment slots for the logged-in admin's hospital
func AdminPage(w http.ResponseWriter, r *http.Request) {
	_, userID, hospID, ok := middleware.GetSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	db, _ := repository.OpenDB()
	defer db.Close()

	var adminName, hospitalName string
	db.QueryRow(`
        SELECT u.full_name, h.name 
        FROM users u 
        JOIN hospital h ON u.hospital_id = h.id 
        WHERE u.id = ?`, userID).Scan(&adminName, &hospitalName)

	var slots []entity.Timeslot
	rows, _ := db.Query(`
        SELECT t.id, d.name, t.doctor, t.room, t.start_time, t.duration,
               t.is_booked, IFNULL(a.patient_name,'')
        FROM timeslot t
        JOIN department d ON t.department_id = d.id
        LEFT JOIN appointment a ON t.id = a.timeslot_id
        WHERE d.hospital_id = ?`, hospID)
	defer rows.Close()

	for rows.Next() {
		var s entity.Timeslot
		var booked int
		rows.Scan(&s.ID, &s.Department, &s.Doctor, &s.Room,
			&s.StartTime, &s.Duration, &booked, &s.Patient)
		s.IsBooked = (booked == 1)
		slots = append(slots, s)
	}

	var depts []entity.Department
	drows, _ := db.Query("SELECT id, name FROM department WHERE hospital_id = ?", hospID)
	defer drows.Close()
	for drows.Next() {
		var d entity.Department
		drows.Scan(&d.ID, &d.Name)
		depts = append(depts, d)
	}

	t := template.Must(template.ParseFiles("templates/admin.html"))
    t.Execute(w, map[string]interface{}{
        "AdminName":    adminName,
        "HospitalName": hospitalName,
        "Timeslots":    slots,
        "Departments":  depts,
    })
}

// AddSlot handles the creation of new timeslots
func AddSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	db, _ := repository.OpenDB()
	defer db.Close()

	db.Exec(`
        INSERT INTO timeslot (department_id, doctor, room, start_time, duration, is_booked)
        VALUES (?, ?, ?, ?, ?, 0)`,
		r.FormValue("department_id"),
		r.FormValue("doctor"),
		r.FormValue("room"),
		r.FormValue("start_time"),
		r.FormValue("duration"),
	)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
