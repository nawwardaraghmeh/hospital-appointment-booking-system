package entity

// user entity, represents either a patient or admin in the system
type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
	FullName     string
	HospitalID   int
}

// city entity
type City struct {
	ID   int
	Name string
}

// hospital entity
type Hospital struct {
	ID     int
	Name   string
	CityID int
}

// department entity
type Department struct {
	ID         int
	Name       string
	HospitalID int
}

// timeslot entity, created by admins, can be booked by patients
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

// appointment entity, signify successfully-booked slots
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

// bookingview entity, to display a patient's existing appointments
type BookingView struct {
	Doctor         string
	StartTime      string
	Room           string
	Symptoms       string
	DepartmentName string
}
