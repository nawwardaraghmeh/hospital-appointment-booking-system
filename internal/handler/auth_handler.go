package handler

import (
	"abs/internal/repository"
	"abs/internal/service"
	"abs/internal/validation"
	"net/http"
	"strconv"
)

// AuthHandler handles registration and login for all user roles
type AuthHandler struct {
	userRepo        repository.UserRepository
	locationRepo    repository.LocationRepository
	adminCode       string
	sessionDuration int
	tmpl            TemplateRenderer
}

func NewAuthHandler(
	userRepo repository.UserRepository,
	locationRepo repository.LocationRepository,
	adminCode string,
	sessionDuration int,
	tmpl TemplateRenderer,
) *AuthHandler {
	return &AuthHandler{
		userRepo:        userRepo,
		locationRepo:    locationRepo,
		adminCode:       adminCode,
		sessionDuration: sessionDuration,
		tmpl:            tmpl,
	}
}

// public home

func (h *AuthHandler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	h.tmpl.Render(w, "index.html", nil)
}

// patient registration

// RegisterPage renders the patient-only public registration form
func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	h.tmpl.Render(w, "register.html", nil)
}

// RegisterAction creates a patient account. role is hardcoded to "patient" here; it is never read from the request
func (h *AuthHandler) RegisterAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	input := &validation.PatientRegisterInput{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		FullName: r.FormValue("full_name"),
	}

	if ve := input.Validate(); ve.HasErrors() {
		h.tmpl.Render(w, "register.html", map[string]interface{}{
			"Errors": ve.Messages,
		})
		return
	}

	hash, err := service.HashPassword(input.Password)
	if err != nil {
		h.tmpl.RenderError(w, "Error processing password.", http.StatusInternalServerError)
		return
	}

	if err := h.userRepo.Create(input.Username, hash, "patient", input.FullName, nil); err != nil {
		h.tmpl.Render(w, "register.html", map[string]interface{}{
			"Errors": []string{"Username already exists."},
		})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// admin registration

// AdminRegisterPage renders the admin registration form
func (h *AuthHandler) AdminRegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	hospitals, err := h.locationRepo.AllHospitals()
	if err != nil {
		h.tmpl.RenderError(w, "Could not load hospitals.", http.StatusInternalServerError)
		return
	}
	h.tmpl.Render(w, "admin_register.html", map[string]interface{}{
		"Hospitals": hospitals,
	})
}

// AdminRegisterAction creates an admin account. Role is hardcoded to "admin" here; it is never read from the request
func (h *AuthHandler) AdminRegisterAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	input := &validation.AdminRegisterInput{
		Username:   r.FormValue("username"),
		Password:   r.FormValue("password"),
		FullName:   r.FormValue("full_name"),
		AdminCode:  r.FormValue("admin_code"),
		HospitalID: r.FormValue("hospital_id"),
	}

	if ve := input.Validate(h.adminCode); ve.HasErrors() {
		hospitals, _ := h.locationRepo.AllHospitals()
		h.tmpl.Render(w, "admin_register.html", map[string]interface{}{
			"Errors":    ve.Messages,
			"Hospitals": hospitals,
		})
		return
	}

	hash, err := service.HashPassword(input.Password)
	if err != nil {
		h.tmpl.RenderError(w, "Error processing password.", http.StatusInternalServerError)
		return
	}

	hospitalID, _ := strconv.Atoi(input.HospitalID)

	if err := h.userRepo.Create(input.Username, hash, "admin", input.FullName, hospitalID); err != nil {
		hospitals, _ := h.locationRepo.AllHospitals()
		h.tmpl.Render(w, "admin_register.html", map[string]interface{}{
			"Errors":    []string{"Username already exists."},
			"Hospitals": hospitals,
		})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// one login

// LoginPage renders the single login form for all users
func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	cities, err := h.locationRepo.AllCities()
	hospitals, err2 := h.locationRepo.AllHospitals()
	departments, err3 := h.locationRepo.AllDepartments()
	if err != nil || err2 != nil || err3 != nil {
		h.tmpl.RenderError(w, "Could not load location data.", http.StatusInternalServerError)
		return
	}
	h.tmpl.Render(w, "login.html", map[string]interface{}{
		"Cities":      cities,
		"Hospitals":   hospitals,
		"Departments": departments,
	})
}

// LoginAction authenticates the user and redirects based on the role resolved from the database
func (h *AuthHandler) LoginAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.tmpl.RenderError(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.userRepo.FindByUsername(r.FormValue("username"))
	if err != nil {
		h.tmpl.RenderError(w, "Invalid username or password.", http.StatusUnauthorized)
		return
	}

	if !service.CheckPasswordHash(r.FormValue("password"), user.PasswordHash) {
		h.tmpl.RenderError(w, "Invalid username or password.", http.StatusUnauthorized)
		return
	}

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

// logout

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
