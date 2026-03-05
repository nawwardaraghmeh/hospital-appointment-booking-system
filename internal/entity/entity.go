package entity

// User represents either a patient or admin in the system
type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
	FullName     string
	HospitalID   int
}

// City entity
type City struct {
	ID   int
	Name string
}

// Hospital entity
type Hospital struct {
	ID     int
	Name   string
	CityID int
}

// Department entity
type Department struct {
	ID         int
	Name       string
	HospitalID int
}

// Timeslot is an available, or already booked, appointment slot created by an admin
type Timeslot struct {
	ID           int
	DepartmentID int
	Department   string
	Doctor       string
	Room         string
	StartTime    string
	Duration     int
	IsBooked     bool
	Patient      string
}

// Appointment represents that successfully booked slots
type Appointment struct {
	ID          int
	TimeSlotID  int
	PatientID   string
	PatientName string
	Age         int
	Phone       string
	Email       string
	Symptoms    string
}

// BookingView is used to display a patient's existing appointments
type BookingView struct {
	Doctor         string
	StartTime      string
	Room           string
	Symptoms       string
	DepartmentName string
}
