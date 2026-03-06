package handler

import (
	"abs/internal/middleware"
	"abs/internal/repository"
	"abs/internal/validation"
	"net/http"
)

// AdminHandler handles all admin dashboard operations
type AdminHandler struct {
	locationRepo repository.LocationRepository
	slotRepo     repository.TimeslotRepository
	tmpl         TemplateRenderer
}

// NewAdminHandler constructs an AdminHandler with all required dependencies injected
func NewAdminHandler(
	locationRepo repository.LocationRepository,
	slotRepo repository.TimeslotRepository,
	tmpl TemplateRenderer,
) *AdminHandler {
	return &AdminHandler{
		locationRepo: locationRepo,
		slotRepo:     slotRepo,
		tmpl:         tmpl,
	}
}

// Dashboard displays all timeslots and the add-slot form for the logged-in admin's hospital
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	_, userID, hospID, ok := middleware.GetSessionCookie(r)
	if !ok {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	adminName, hospitalName, err := h.locationRepo.FindHospitalNameByUserID(userID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load admin info.", http.StatusInternalServerError)
		return
	}

	slots, err := h.slotRepo.FindAllByHospital(hospID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load timeslots.", http.StatusInternalServerError)
		return
	}

	depts, err := h.locationRepo.DepartmentsByHospital(hospID)
	if err != nil {
		h.tmpl.RenderError(w, "Could not load departments.", http.StatusInternalServerError)
		return
	}

	h.tmpl.Render(w, "admin.html", map[string]interface{}{
		"AdminName":    adminName,
		"HospitalName": hospitalName,
		"Timeslots":    slots,
		"Departments":  depts,
	})
}

// AddSlot processes the form submission that creates a new timeslot
func (h *AdminHandler) AddSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	input := &validation.TimeslotInput{
		DepartmentID: r.FormValue("department_id"),
		Doctor:       r.FormValue("doctor"),
		Room:         r.FormValue("room"),
		StartTime:    r.FormValue("start_time"),
		DurationStr:  r.FormValue("duration"),
	}

	duration, ve := input.Validate()
	if ve.HasErrors() {
		_, userID, hospID, ok := middleware.GetSessionCookie(r)
		if !ok {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		adminName, hospitalName, _ := h.locationRepo.FindHospitalNameByUserID(userID)
		slots, _ := h.slotRepo.FindAllByHospital(hospID)
		depts, _ := h.locationRepo.DepartmentsByHospital(hospID)

		h.tmpl.Render(w, "admin.html", map[string]interface{}{
			"AdminName":    adminName,
			"HospitalName": hospitalName,
			"Timeslots":    slots,
			"Departments":  depts,
			"Errors":       ve.Messages,
		})
		return
	}

	if err := h.slotRepo.Create(input.DepartmentID, input.Doctor, input.Room, input.StartTime, duration); err != nil {
		h.tmpl.RenderError(w, "Could not create timeslot.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}
