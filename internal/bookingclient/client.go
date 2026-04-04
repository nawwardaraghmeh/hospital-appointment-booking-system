package bookingclient

import (
	"abs/internal/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client communicates with the Booking Service REST API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a Client that calls the Booking Service
func New(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// AllCities returns all cities from the Booking Service
func (c *Client) AllCities() ([]entity.City, error) {
	var result []entity.City
	if err := c.get("/api/cities", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AllHospitals returns all hospitals from the Booking Service
func (c *Client) AllHospitals() ([]entity.Hospital, error) {
	var result []entity.Hospital
	if err := c.get("/api/hospitals", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AllDepartments returns all departments from the Booking Service
func (c *Client) AllDepartments() ([]entity.Department, error) {
	var result []entity.Department
	if err := c.get("/api/departments", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// DepartmentsByHospital returns departments for a specific hospital
func (c *Client) DepartmentsByHospital(hospitalID int) ([]entity.Department, error) {
	var result []entity.Department
	if err := c.get("/api/departments?hospital_id="+strconv.Itoa(hospitalID), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// FindDepartmentWithHospital returns the department and hospital name for a dept ID
func (c *Client) FindDepartmentWithHospital(deptID string) (deptName, hospName string, err error) {
	var result struct {
		DepartmentName string `json:"department_name"`
		HospitalName   string `json:"hospital_name"`
	}
	if err := c.get("/api/departments/"+deptID+"/hospital", &result); err != nil {
		return "", "", err
	}
	return result.DepartmentName, result.HospitalName, nil
}

// FindHospitalNameByUserID returns the admin name and hospital name for a given user ID
func (c *Client) FindHospitalNameByUserID(userID int) (adminName, hospitalName string, err error) {
	var result struct {
		AdminName    string `json:"admin_name"`
		HospitalName string `json:"hospital_name"`
	}
	if err := c.get("/api/admin/"+strconv.Itoa(userID)+"/hospital", &result); err != nil {
		return "", "", err
	}
	return result.AdminName, result.HospitalName, nil
}

// FindAvailableByDepartment returns unbooked slots for a department
func (c *Client) FindAvailableByDepartment(deptID string) ([]entity.Timeslot, error) {
	var result []entity.Timeslot
	if err := c.get("/api/slots?department_id="+deptID+"&available=true", &result); err != nil {
		return nil, err
	}
	return result, nil
}

// FindAllByHospital returns all slots for a hospital (admin view)
func (c *Client) FindAllByHospital(hospitalID int) ([]entity.Timeslot, error) {
	var result []entity.Timeslot
	if err := c.get("/api/slots?hospital_id="+strconv.Itoa(hospitalID), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// FindDepartmentIDBySlotID returns the department ID string for a slot
func (c *Client) FindDepartmentIDBySlotID(slotID int) (string, error) {
	var result struct {
		DepartmentID string `json:"department_id"`
	}
	if err := c.get("/api/slots/"+strconv.Itoa(slotID)+"/department", &result); err != nil {
		return "", err
	}
	return result.DepartmentID, nil
}

// CreateSlot asks the Booking Service to create a new timeslot
func (c *Client) CreateSlot(departmentID, doctor, room, startTime string, duration int) error {
	body := fmt.Sprintf(
		`{"department_id":%q,"doctor":%q,"room":%q,"start_time":%q,"duration":%d}`,
		departmentID, doctor, room, startTime, duration,
	)
	return c.post("/api/slots", body)
}

// MyAppointments returns a patient's booked appointments
func (c *Client) MyAppointments(patientID int) ([]entity.BookingView, error) {
	var result []entity.BookingView
	if err := c.get("/api/appointments?patient_id="+strconv.Itoa(patientID), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Reserve sends a booking request to the Booking Service
func (c *Client) Reserve(appt entity.Appointment) error {
	body := fmt.Sprintf(
		`{"timeslot_id":%d,"patient_id":%q,"patient_name":%q,"age":%d,"phone":%q,"email":%q,"symptoms":%q}`,
		appt.TimeSlotID, appt.PatientID, appt.PatientName, appt.Age, appt.Phone, appt.Email, appt.Symptoms,
	)
	return c.post("/api/appointments", body)
}

// HTTP helpers
func (c *Client) get(path string, out interface{}) error {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("booking service unreachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("booking service returned %d for GET %s", resp.StatusCode, path)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) post(path, jsonBody string) error {
	resp, err := c.httpClient.Post(
		c.baseURL+path,
		"application/json",
		strings.NewReader(jsonBody),
	)
	if err != nil {
		return fmt.Errorf("booking service unreachable: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("booking service error %d for POST %s", resp.StatusCode, path)
	}
	return nil
}
