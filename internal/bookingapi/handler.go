package bookingapi

import (
	"abs/internal/entity"
	"abs/internal/repository"
	"abs/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Handler handles all Booking Service API routes
type Handler struct {
	locationRepo repository.LocationRepository
	slotRepo     repository.TimeslotRepository
	bookingSvc   service.BookingService
}

// New constructs a Handler with all required dependencies
func New(
	locationRepo repository.LocationRepository,
	slotRepo repository.TimeslotRepository,
	bookingSvc service.BookingService,
) *Handler {
	return &Handler{
		locationRepo: locationRepo,
		slotRepo:     slotRepo,
		bookingSvc:   bookingSvc,
	}
}

// RegisterRoutes wires all API routes onto mux
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/cities", h.cities)
	mux.HandleFunc("/api/hospitals", h.hospitals)
	mux.HandleFunc("/api/departments", h.departments)
	mux.HandleFunc("/api/departments/", h.departmentDetail)
	mux.HandleFunc("/api/admin/", h.adminHospital)
	mux.HandleFunc("/api/slots", h.slots)
	mux.HandleFunc("/api/slots/", h.slotDetail)
	mux.HandleFunc("/api/appointments", h.appointments)
}

// location endpoints
func (h *Handler) cities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cities, err := h.locationRepo.AllCities()
	if err != nil {
		writeError(w, "could not load cities", http.StatusInternalServerError)
		return
	}
	writeJSON(w, cities)
}

func (h *Handler) hospitals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	hospitals, err := h.locationRepo.AllHospitals()
	if err != nil {
		writeError(w, "could not load hospitals", http.StatusInternalServerError)
		return
	}
	writeJSON(w, hospitals)
}

func (h *Handler) departments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if hospID := r.URL.Query().Get("hospital_id"); hospID != "" {
		depts, err := h.locationRepo.DepartmentsByHospital(hospID)
		if err != nil {
			writeError(w, "could not load departments", http.StatusInternalServerError)
			return
		}
		writeJSON(w, depts)
		return
	}

	depts, err := h.locationRepo.AllDepartments()
	if err != nil {
		writeError(w, "could not load departments", http.StatusInternalServerError)
		return
	}
	writeJSON(w, depts)
}

// departmentDetail handles GET /api/departments/{id}/hospital
func (h *Handler) departmentDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Path: /api/departments/{id}/hospital → parts[2] = id
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, "invalid path", http.StatusBadRequest)
		return
	}
	deptID := parts[2]

	deptName, hospName, err := h.locationRepo.FindDepartmentWithHospital(deptID)
	if err != nil {
		writeError(w, "department not found", http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]string{
		"department_name": deptName,
		"hospital_name":   hospName,
	})
}

// adminHospital handles GET /api/admin/{userID}/hospital
func (h *Handler) adminHospital(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, "invalid path", http.StatusBadRequest)
		return
	}
	userID := parts[2]

	adminName, hospitalName, err := h.locationRepo.FindHospitalNameByUserID(userID)
	if err != nil {
		writeError(w, "admin not found", http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]string{
		"admin_name":    adminName,
		"hospital_name": hospitalName,
	})
}

// slot endpoints
func (h *Handler) slots(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getSlots(w, r)
	case http.MethodPost:
		h.createSlot(w, r)
	default:
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getSlots(w http.ResponseWriter, r *http.Request) {
	if deptID := r.URL.Query().Get("department_id"); deptID != "" {
		slots, err := h.slotRepo.FindAvailableByDepartment(deptID)
		if err != nil {
			writeError(w, "could not load slots", http.StatusInternalServerError)
			return
		}
		writeJSON(w, slots)
		return
	}

	if hospID := r.URL.Query().Get("hospital_id"); hospID != "" {
		slots, err := h.slotRepo.FindAllByHospital(hospID)
		if err != nil {
			writeError(w, "could not load slots", http.StatusInternalServerError)
			return
		}
		writeJSON(w, slots)
		return
	}

	writeError(w, "department_id or hospital_id query param required", http.StatusBadRequest)
}

func (h *Handler) createSlot(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DepartmentID string `json:"department_id"`
		Doctor       string `json:"doctor"`
		Room         string `json:"room"`
		StartTime    string `json:"start_time"`
		Duration     int    `json:"duration"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.slotRepo.Create(body.DepartmentID, body.Doctor, body.Room, body.StartTime, body.Duration); err != nil {
		writeError(w, "could not create slot", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// slotDetail handles GET /api/slots/{id}/department
func (h *Handler) slotDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, "invalid path", http.StatusBadRequest)
		return
	}
	slotIDInt := 0
	fmt.Sscanf(parts[2], "%d", &slotIDInt)

	deptID, err := h.slotRepo.FindDepartmentIDBySlotID(slotIDInt)
	if err != nil {
		writeError(w, "slot not found", http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]string{"department_id": deptID})
}

// appointment endpoints
func (h *Handler) appointments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAppointments(w, r)
	case http.MethodPost:
		h.createAppointment(w, r)
	default:
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getAppointments(w http.ResponseWriter, r *http.Request) {
	patientID := r.URL.Query().Get("patient_id")
	appts, err := h.bookingSvc.MyAppointments(patientID)
	if err != nil {
		writeError(w, "could not load appointments", http.StatusInternalServerError)
		return
	}
	writeJSON(w, appts)
}

func (h *Handler) createAppointment(w http.ResponseWriter, r *http.Request) {
	var body struct {
		TimeSlotID  int    `json:"timeslot_id"`
		PatientID   string `json:"patient_id"`
		PatientName string `json:"patient_name"`
		Age         int    `json:"age"`
		Phone       string `json:"phone"`
		Email       string `json:"email"`
		Symptoms    string `json:"symptoms"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	appt := entity.Appointment{
		TimeSlotID:  body.TimeSlotID,
		PatientID:   body.PatientID,
		PatientName: body.PatientName,
		Age:         body.Age,
		Phone:       body.Phone,
		Email:       body.Email,
		Symptoms:    body.Symptoms,
	}

	if err := h.bookingSvc.Reserve(appt); err != nil {
		writeError(w, "booking failed: slot may already be taken", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// helpers
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
