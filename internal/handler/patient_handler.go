package handler

import (
	"abs/internal/entity"
	"abs/internal/middleware"
	"abs/internal/repository"
	"abs/internal/service"
	"html/template"
	"net/http"
	"strconv"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	t.Execute(w, nil)
}

func LogoutPage(w http.ResponseWriter, r *http.Request) {
	middleware.ClearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// PatientLoginAction validates the login form and establishes a patient session
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

// PatientLoginPage loads all data for the dynamic population of dropdown menus
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

// PatientSlotsPage shows 1. available slots, and 2. the patient's existing bookings.
func PatientSlotsPage(w http.ResponseWriter, r *http.Request) {
	deptID := r.URL.Query().Get("department_id")
	if deptID == "" {
		http.Error(w, "No department selected", 400)
		return
	}

	_, patientID, ok := middleware.GetSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/patient/login", http.StatusSeeOther)
		return
	}

	db, err := repository.OpenDB()
	if err != nil {
		http.Error(w, "Database connection error", 500)
		return
	}
	defer db.Close()

	// 1. fetch unbooked slots for the chosen department
	var slots []entity.Timeslot
	rows, _ := db.Query(`
        SELECT id, doctor, room, start_time, duration
        FROM timeslot
        WHERE department_id=? AND is_booked=0
    `, deptID)
	defer rows.Close()

	for rows.Next() {
		var s entity.Timeslot
		rows.Scan(&s.ID, &s.Doctor, &s.Room, &s.StartTime, &s.Duration)
		slots = append(slots, s)
	}

	// 2. fetch the appointments the logged-in patient booked
	var myAppointments []entity.BookingView
	appRows, _ := db.Query(`
        SELECT t.doctor, t.start_time, t.room, a.symptoms, d.name
        FROM appointment a
        JOIN timeslot t ON a.timeslot_id = t.id
        JOIN department d ON t.department_id = d.id
        WHERE a.patient_id = ?`, patientID)
	defer appRows.Close()

	for appRows.Next() {
		var bv entity.BookingView
		err := appRows.Scan(&bv.Doctor, &bv.StartTime, &bv.Room, &bv.Symptoms, &bv.DepartmentName)
		if err != nil {
			continue
		}
		myAppointments = append(myAppointments, bv)
	}

	bookedSuccess := r.URL.Query().Get("booked") == "true"

	t := template.Must(template.ParseFiles("templates/patient_slots.html"))
	t.Execute(w, map[string]interface{}{
		"Slots":          slots,
		"MyAppointments": myAppointments,
		"DepartmentID":   deptID,
		"BookedSuccess":  bookedSuccess,
	})
}

// BookAppointment processes the booking form and updates the slot status
func BookAppointment(w http.ResponseWriter, r *http.Request) {
	slotIDStr := r.URL.Query().Get("slot_id")

	slotID, err := strconv.Atoi(slotIDStr)
	if err != nil {
		http.Error(w, "Invalid Slot ID", http.StatusBadRequest)
		return
	}

	db, _ := repository.OpenDB()
	defer db.Close()

	if r.Method == http.MethodPost {
		age, _ := strconv.Atoi(r.FormValue("age"))

		appt := entity.Appointment{
			TimeSlotID:  slotID,
			PatientID:   r.FormValue("pid"),
			PatientName: r.FormValue("name"),
			Age:         age,
			Phone:       r.FormValue("phone"),
			Email:       r.FormValue("email"),
			Symptoms:    r.FormValue("symptoms"),
		}

		err := service.Reserve(db, appt)
		if err != nil {
			http.Error(w, "Booking failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var deptID string
		db.QueryRow("SELECT department_id FROM timeslot WHERE id = ?", slotID).Scan(&deptID)

		http.Redirect(w, r, "/patient/slots?department_id="+deptID+"&booked=true", http.StatusSeeOther)
		return
	}

	t := template.Must(template.ParseFiles("templates/book.html"))
	t.Execute(w, map[string]interface{}{
		"SlotID": slotID,
	})
}
