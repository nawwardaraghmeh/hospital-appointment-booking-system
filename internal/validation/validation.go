package validation

import (
	"fmt"
	"strings"
)

// ValidationError holds a list of user-facing error messages
type ValidationError struct {
	Messages []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Messages, "; ")
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Messages) > 0
}

func (e *ValidationError) Add(msg string) {
	e.Messages = append(e.Messages, msg)
}

// PatientRegisterInput validates the public registration form
type PatientRegisterInput struct {
	Username string
	Password string
	FullName string
}

func (i *PatientRegisterInput) Validate() *ValidationError {
	ve := &ValidationError{}
	if strings.TrimSpace(i.Username) == "" {
		ve.Add("Username is required.")
	}
	if len(i.Password) < 6 {
		ve.Add("Password must be at least 6 characters.")
	}
	if strings.TrimSpace(i.FullName) == "" {
		ve.Add("Full name is required.")
	}
	return ve
}

// AdminRegisterInput validates the hidden admin registration form
type AdminRegisterInput struct {
	Username   string
	Password   string
	FullName   string
	AdminCode  string
	HospitalID string
}

func (i *AdminRegisterInput) Validate(expectedAdminCode string) *ValidationError {
	ve := &ValidationError{}
	if strings.TrimSpace(i.Username) == "" {
		ve.Add("Username is required.")
	}
	if len(i.Password) < 6 {
		ve.Add("Password must be at least 6 characters.")
	}
	if strings.TrimSpace(i.FullName) == "" {
		ve.Add("Full name is required.")
	}
	if i.AdminCode != expectedAdminCode {
		ve.Add("Invalid admin registration code.")
	}
	if strings.TrimSpace(i.HospitalID) == "" || i.HospitalID == "0" {
		ve.Add("You must select a hospital.")
	}
	return ve
}

// BookingInput holds and validates the appointment booking form fields
type BookingInput struct {
	Name     string
	AgeStr   string
	Phone    string
	Email    string
	Symptoms string
	SlotID   int
}

func (i *BookingInput) Validate() (*BookingInput, *ValidationError) {
	ve := &ValidationError{}
	if strings.TrimSpace(i.Name) == "" {
		ve.Add("Full name is required.")
	}
	age := 0
	if strings.TrimSpace(i.AgeStr) == "" {
		ve.Add("Age is required.")
	} else {
		n, err := fmt.Sscanf(i.AgeStr, "%d", &age)
		if err != nil || n == 0 || age < 0 || age > 150 {
			ve.Add("Age must be a number between 0 and 150.")
		}
	}
	if i.SlotID <= 0 {
		ve.Add("Invalid slot selected.")
	}
	return i, ve
}

// TimeslotInput holds and validates the add-slot form fields
type TimeslotInput struct {
	DepartmentID string
	Doctor       string
	Room         string
	StartTime    string
	DurationStr  string
}

func (i *TimeslotInput) Validate() (duration int, ve *ValidationError) {
	ve = &ValidationError{}
	if strings.TrimSpace(i.DepartmentID) == "" || i.DepartmentID == "0" {
		ve.Add("Department is required.")
	}
	if strings.TrimSpace(i.Doctor) == "" {
		ve.Add("Doctor name is required.")
	}
	if strings.TrimSpace(i.Room) == "" {
		ve.Add("Room is required.")
	}
	if strings.TrimSpace(i.StartTime) == "" {
		ve.Add("Start time is required.")
	}
	n, err := fmt.Sscanf(i.DurationStr, "%d", &duration)
	if err != nil || n == 0 || duration <= 0 || duration > 480 {
		ve.Add("Duration must be a positive number of minutes (max 480).")
	}
	return duration, ve
}
