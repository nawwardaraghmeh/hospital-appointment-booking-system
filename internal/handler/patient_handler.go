package handler

import (
	"abs/internal/entity"
	"abs/internal/middleware"
	"abs/internal/repository"
	"abs/internal/service"
	"abs/internal/validation"
	"net/http"
	"strconv"
)

// PatientHandler handles all patient operations
type PatientHandler struct {
	userRepo     repository.UserRepository
	locationRepo repository.LocationRepository
	slotRepo     repository.TimeslotRepository
	bookingSvc   service.BookingService
	tmpl         TemplateRenderer
}

// NewPatientHandler constructs a PatientHandler with all required dependencies injected
func NewPatientHandler(
	userRepo repository.UserRepository,
	locationRepo repository.LocationRepository,
	slotRepo repository.TimeslotRepository,
	bookingSvc service.BookingService,
	tmpl TemplateRenderer,
) *PatientHandler {
	return &PatientHandler{
		userRepo:     userRepo,
		locationRepo: locationRepo,
		slotRepo:     slotRepo,
		bookingSvc:   bookingSvc,
		tmpl:         tmpl,
	}
}

// SlotsPage shows available slots for a department and the patient's existing bookings
func (h *PatientHandler) SlotsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	deptID := r.URL.Query().Get("department_id")
	if deptID == "" {
		h.tmpl.RenderError(w, "No department selected.", http.StatusBadRequest)
		return
	}

	_, patientID, _, ok := middleware.GetSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/patient/login", http.StatusSeeOther)
		return
	}

	patient, err := h.userRepo.FindByID(patientID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load patient info.", http.StatusInternalServerError)
		return
	}

	deptName, hospName, err := h.locationRepo.FindDepartmentWithHospital(deptID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load department info.", http.StatusInternalServerError)
		return
	}

	slots, err := h.slotRepo.FindAvailableByDepartment(deptID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load available slots.", http.StatusInternalServerError)
		return
	}

	myAppointments, err := h.bookingSvc.MyAppointments(patientID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load your appointments.", http.StatusInternalServerError)
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

// BookPage renders the booking form and processes the submission
func (h *PatientHandler) BookPage(w http.ResponseWriter, r *http.Request) {
	slotIDStr := r.URL.Query().Get("slot_id")
	slotID, _ := strconv.Atoi(slotIDStr)

	_, userID, _, ok := middleware.GetSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		h.tmpl.Render(w, "book.html", map[string]interface{}{"SlotID": slotID})
		return
	}

	if r.Method != http.MethodPost {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	input := &validation.BookingInput{
		Name:     r.FormValue("name"),
		AgeStr:   r.FormValue("age"),
		Phone:    r.FormValue("phone"),
		Email:    r.FormValue("email"),
		Symptoms: r.FormValue("symptoms"),
		SlotID:   slotID,
	}

	_, ve := input.Validate()
	if ve.HasErrors() {
		h.tmpl.Render(w, "book.html", map[string]interface{}{
			"SlotID":   slotID,
			"Errors":   ve.Messages,
			"Name":     input.Name,
			"Age":      input.AgeStr,
			"Phone":    input.Phone,
			"Email":    input.Email,
			"Symptoms": input.Symptoms,
		})
		return
	}

	age, _ := strconv.Atoi(input.AgeStr)
	appt := entity.Appointment{
		TimeSlotID:  slotID,
		PatientID:   userID,
		PatientName: input.Name,
		Age:         age,
		Phone:       input.Phone,
		Email:       input.Email,
		Symptoms:    input.Symptoms,
	}

	if err := h.bookingSvc.Reserve(appt); err != nil {
		h.tmpl.RenderError(w, "Booking failed — this slot may already be taken.", http.StatusConflict)
		return
	}

	deptID, err := h.slotRepo.FindDepartmentIDBySlotID(slotID)
	if err != nil {
		http.Redirect(w, r, "/patient/slots", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/patient/slots?department_id="+deptID+"&booked=true", http.StatusSeeOther)
}
