package handler

import (
	"abs/internal/entity"
	"abs/internal/repository"
	"html/template"
	"net/http"
)

func AdminLoginPage(w http.ResponseWriter, r *http.Request) {
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

	t := template.Must(template.ParseFiles("templates/admin_login.html"))
	t.Execute(w, map[string]interface{}{
		"Cities":    cities,
		"Hospitals": hospitals,
	})
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	db, _ := repository.OpenDB()
	defer db.Close()

	var slots []entity.Timeslot
	rows, _ := db.Query(`
		SELECT t.id, d.name, t.doctor, t.room, t.start_time, t.duration,
		       t.is_booked, IFNULL(a.patient_name,'')
		FROM timeslot t
		JOIN department d ON t.department_id=d.id
		LEFT JOIN appointment a ON t.id=a.timeslot_id
	`)
	defer rows.Close()

	for rows.Next() {
		var s entity.Timeslot
		var booked int
		rows.Scan(
			&s.ID, &s.Department, &s.Doctor, &s.Room,
			&s.StartTime, &s.Duration, &booked, &s.Patient,
		)
		s.IsBooked = booked == 1
		slots = append(slots, s)
	}

	var depts []entity.Department
	drows, _ := db.Query("SELECT id, name, hospital_id FROM department")
	defer drows.Close()
	for drows.Next() {
		var d entity.Department
		drows.Scan(&d.ID, &d.Name, &d.HospitalID)
		depts = append(depts, d)
	}

	t := template.Must(template.ParseFiles("templates/admin.html"))
	t.Execute(w, map[string]interface{}{
		"Timeslots":   slots,
		"Departments": depts,
	})
}

func AddSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		AdminPage(w, r)
		return
	}

	db, _ := repository.OpenDB()
	defer db.Close()

	db.Exec(`
		INSERT INTO timeslot
		(department_id, doctor, room, start_time, duration, is_booked)
		VALUES (?, ?, ?, ?, ?, 0)`,
		r.FormValue("department_id"),
		r.FormValue("doctor"),
		r.FormValue("room"),
		r.FormValue("start_time"),
		r.FormValue("duration"),
	)

	AdminPage(w, r)
}
