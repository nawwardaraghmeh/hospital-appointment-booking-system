package handler

import (
	"abs/internal/bookingclient"
	"abs/internal/entity"
	"abs/internal/middleware"
	"abs/internal/repository"
	"abs/internal/service"
	"abs/internal/validation"
	"net/http"
	"strconv"
)

// handle the dependencies: registration, login, logout, and all patient/admin pages
type AuthHandler struct {
	userRepo        repository.UserRepository
	bookingClient   *bookingclient.Client
	adminCode       string
	sessionDuration int
	tmpl            TemplateRenderer
}

// create an AuthHandler with all required dependencies
func NewAuthHandler(
	userRepo repository.UserRepository,
	bookingClient *bookingclient.Client,
	adminCode string,
	sessionDuration int,
	tmpl TemplateRenderer,
) *AuthHandler {
	return &AuthHandler{
		userRepo:        userRepo,
		bookingClient:   bookingClient,
		adminCode:       adminCode,
		sessionDuration: sessionDuration,
		tmpl:            tmpl,
	}
}

// the main landing page
func (h *AuthHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	h.tmpl.Render(w, "index.html", nil)
}

// the patient registration form
func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.tmpl.Render(w, "register.html", nil)
}

// process and validate patient registration
func (h *AuthHandler) RegisterAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.tmpl.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	registerInput := &validation.PatientRegisterInput{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		FullName: r.FormValue("full_name"),
	}

	if ve := registerInput.Validate(); ve.HasErrors() {
		h.tmpl.Render(w, "register.html", map[string]interface{}{
			"Errors":   ve.Messages,
			"Username": registerInput.Username,
			"FullName": registerInput.FullName,
		})
		return
	}

	hash, err := service.HashPassword(registerInput.Password)
	if err != nil {
		h.tmpl.RenderError(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	if err := h.userRepo.Create(registerInput.Username, hash, "patient", registerInput.FullName, nil); err != nil {
		h.tmpl.Render(w, "register.html", map[string]interface{}{
			"Errors":   []string{"Username already exists"},
			"Username": registerInput.Username,
			"FullName": registerInput.FullName,
		})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// the hidden admin registration form
func (h *AuthHandler) AdminRegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hospitals, err := h.bookingClient.AllHospitals()
	if err != nil {
		h.tmpl.RenderError(w, "Could not load hospitals from Booking Service", http.StatusBadGateway)
		return
	}
	h.tmpl.Render(w, "admin_register.html", map[string]interface{}{
		"Hospitals": hospitals,
	})
}

// process and validate admin registration
func (h *AuthHandler) AdminRegisterAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.tmpl.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	adminInput := &validation.AdminRegisterInput{
		Username:   r.FormValue("username"),
		Password:   r.FormValue("password"),
		FullName:   r.FormValue("full_name"),
		AdminCode:  r.FormValue("admin_code"),
		HospitalID: r.FormValue("hospital_id"),
	}

	if ve := adminInput.Validate(h.adminCode); ve.HasErrors() {
		hospitals, _ := h.bookingClient.AllHospitals()
		h.tmpl.Render(w, "admin_register.html", map[string]interface{}{
			"Hospitals":  hospitals,
			"Errors":     ve.Messages,
			"Username":   adminInput.Username,
			"FullName":   adminInput.FullName,
			"HospitalID": adminInput.HospitalID,
		})
		return
	}

	hash, err := service.HashPassword(adminInput.Password)
	if err != nil {
		h.tmpl.RenderError(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	hospitalID, _ := strconv.Atoi(adminInput.HospitalID)
	if err := h.userRepo.Create(adminInput.Username, hash, "admin", adminInput.FullName, hospitalID); err != nil {
		hospitals, _ := h.bookingClient.AllHospitals()
		h.tmpl.Render(w, "admin_register.html", map[string]interface{}{
			"Hospitals":  hospitals,
			"Errors":     []string{"Username already exists"},
			"Username":   adminInput.Username,
			"FullName":   adminInput.FullName,
			"HospitalID": adminInput.HospitalID,
		})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// the shared login form
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cities, err := h.bookingClient.AllCities()
	if err != nil {
		h.tmpl.RenderError(w, "Could not load cities from Booking Service", http.StatusBadGateway)
		return
	}
	hospitals, err := h.bookingClient.AllHospitals()
	if err != nil {
		h.tmpl.RenderError(w, "Could not load hospitals from Booking Service", http.StatusBadGateway)
		return
	}
	departments, err := h.bookingClient.AllDepartments()
	if err != nil {
		h.tmpl.RenderError(w, "Could not load departments from Booking Service", http.StatusBadGateway)
		return
	}

	h.tmpl.Render(w, "login.html", map[string]interface{}{
		"Cities":      cities,
		"Hospitals":   hospitals,
		"Departments": departments,
	})
}

// authenticate the user and set the session cookie
func (h *AuthHandler) LoginAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.tmpl.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.userRepo.FindByUsername(r.FormValue("username"))
	if err != nil {
		h.tmpl.Render(w, "login.html", map[string]interface{}{"Error": "Invalid username or password"})
		return
	}

	if !service.CheckPasswordHash(r.FormValue("password"), user.PasswordHash) {
		h.tmpl.Render(w, "login.html", map[string]interface{}{"Error": "Invalid username or password"})
		return
	}

	// role read from DB
	cookieValue := user.Role + ":" + strconv.Itoa(user.ID) + ":" + strconv.Itoa(user.HospitalID)
	http.SetCookie(w, &http.Cookie{
		Name:     "abs_session",
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   h.sessionDuration,
	})

	if user.Role == "admin" {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	if deptID := r.FormValue("department_id"); deptID != "" {
		http.Redirect(w, r, "/patient/slots?department_id="+deptID, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/patient/slots", http.StatusSeeOther)
	}
}

// destroy the session cookie upon logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "abs_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// the admin dashboard
func (h *AuthHandler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, userID, hospID, ok := middleware.GetSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	adminName, hospitalName, err := h.bookingClient.FindHospitalNameByUserID(userID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load admin info from Booking Service", http.StatusBadGateway)
		return
	}

	slots, err := h.bookingClient.FindAllByHospital(hospID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load timeslots from Booking Service", http.StatusBadGateway)
		return
	}

	depts, err := h.bookingClient.DepartmentsByHospital(hospID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load departments from Booking Service", http.StatusBadGateway)
		return
	}

	h.tmpl.Render(w, "admin.html", map[string]interface{}{
		"AdminName":    adminName,
		"HospitalName": hospitalName,
		"Timeslots":    slots,
		"Departments":  depts,
	})
}

// process and validate a new timeslot before sending it to the booking service
func (h *AuthHandler) AddSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.tmpl.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, userID, hospID, ok := middleware.GetSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	timeslotInput := &validation.TimeslotInput{
		DepartmentID: r.FormValue("department_id"),
		Doctor:       r.FormValue("doctor"),
		Room:         r.FormValue("room"),
		StartTime:    r.FormValue("start_time"),
		DurationStr:  r.FormValue("duration"),
	}

	duration, ve := timeslotInput.Validate()
	if ve.HasErrors() {
		adminName, hospitalName, _ := h.bookingClient.FindHospitalNameByUserID(userID)
		slots, _ := h.bookingClient.FindAllByHospital(hospID)
		depts, _ := h.bookingClient.DepartmentsByHospital(hospID)
		h.tmpl.Render(w, "admin.html", map[string]interface{}{
			"AdminName":    adminName,
			"HospitalName": hospitalName,
			"Timeslots":    slots,
			"Departments":  depts,
			"Errors":       ve.Messages,
			"DepartmentID": timeslotInput.DepartmentID,
			"Doctor":       timeslotInput.Doctor,
			"Room":         timeslotInput.Room,
			"StartTime":    timeslotInput.StartTime,
			"Duration":     timeslotInput.DurationStr,
		})
		return
	}

	if err := h.bookingClient.CreateSlot(
		timeslotInput.DepartmentID,
		timeslotInput.Doctor,
		timeslotInput.Room,
		timeslotInput.StartTime,
		duration,
	); err != nil {
		adminName, hospitalName, _ := h.bookingClient.FindHospitalNameByUserID(userID)
		slots, _ := h.bookingClient.FindAllByHospital(hospID)
		depts, _ := h.bookingClient.DepartmentsByHospital(hospID)
		h.tmpl.Render(w, "admin.html", map[string]interface{}{
			"AdminName":    adminName,
			"HospitalName": hospitalName,
			"Timeslots":    slots,
			"Departments":  depts,
			"Errors":       []string{"Could not create timeslot: " + err.Error()},
			"DepartmentID": timeslotInput.DepartmentID,
			"Doctor":       timeslotInput.Doctor,
			"Room":         timeslotInput.Room,
			"StartTime":    timeslotInput.StartTime,
			"Duration":     timeslotInput.DurationStr,
		})
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// fetche available slots
func (h *AuthHandler) SlotsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deptID := r.URL.Query().Get("department_id")
	if deptID == "" {
		h.tmpl.RenderError(w, "No department selected", http.StatusBadRequest)
		return
	}

	_, patientID, _, ok := middleware.GetSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	patient, err := h.userRepo.FindByID(strconv.Itoa(patientID))
	if err != nil {
		h.tmpl.RenderError(w, "Could not load patient info", http.StatusInternalServerError)
		return
	}

	deptName, hospName, err := h.bookingClient.FindDepartmentWithHospital(deptID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load department info from Booking Service", http.StatusBadGateway)
		return
	}

	slots, err := h.bookingClient.FindAvailableByDepartment(deptID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load slots from Booking Service", http.StatusBadGateway)
		return
	}

	myAppointments, err := h.bookingClient.MyAppointments(patientID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load appointments from Booking Service", http.StatusBadGateway)
		return
	}

	h.tmpl.Render(w, "patient_slots.html", map[string]interface{}{
		"PatientName":    patient.FullName,
		"HospitalName":   hospName,
		"DepartmentName": deptName,
		"Slots":          slots,
		"MyAppointments": myAppointments,
		"DepartmentID":   deptID,
		"BookedSuccess":  r.URL.Query().Get("booked") == "true",
	})
}

// the booking form (GET). and validate and submit the booking (POST)
func (h *AuthHandler) BookPage(w http.ResponseWriter, r *http.Request) {
	slotIDStr := r.URL.Query().Get("slot_id")
	slotID, _ := strconv.Atoi(slotIDStr)

	_, userID, _, ok := middleware.GetSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		bookingInput := &validation.BookingInput{
			Name:     r.FormValue("name"),
			AgeStr:   r.FormValue("age"),
			Phone:    r.FormValue("phone"),
			Email:    r.FormValue("email"),
			Symptoms: r.FormValue("symptoms"),
			SlotID:   slotID,
		}

		_, ve := bookingInput.Validate()
		if ve.HasErrors() {
			h.tmpl.Render(w, "book.html", map[string]interface{}{
				"SlotID":   slotID,
				"Errors":   ve.Messages,
				"Name":     bookingInput.Name,
				"Age":      bookingInput.AgeStr,
				"Phone":    bookingInput.Phone,
				"Email":    bookingInput.Email,
				"Symptoms": bookingInput.Symptoms,
			})
			return
		}

		age, _ := strconv.Atoi(bookingInput.AgeStr)

		appt := entity.Appointment{
			TimeSlotID:  slotID,
			PatientID:   strconv.Itoa(userID),
			PatientName: bookingInput.Name,
			Age:         age,
			Phone:       bookingInput.Phone,
			Email:       bookingInput.Email,
			Symptoms:    bookingInput.Symptoms,
		}

		if err := h.bookingClient.Reserve(appt); err != nil {
			h.tmpl.Render(w, "book.html", map[string]interface{}{
				"SlotID":   slotID,
				"Errors":   []string{"Booking failed: slot may already be taken"},
				"Name":     bookingInput.Name,
				"Age":      bookingInput.AgeStr,
				"Phone":    bookingInput.Phone,
				"Email":    bookingInput.Email,
				"Symptoms": bookingInput.Symptoms,
			})
			return
		}

		deptID, err := h.bookingClient.FindDepartmentIDBySlotID(slotID)
		if err != nil {
			http.Redirect(w, r, "/patient/slots", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/patient/slots?department_id="+deptID+"&booked=true", http.StatusSeeOther)
		return
	}

	h.tmpl.Render(w, "book.html", map[string]interface{}{"SlotID": slotID})
}
