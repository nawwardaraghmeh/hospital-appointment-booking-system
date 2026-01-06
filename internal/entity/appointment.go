package entity

// this entity represents the data stored when a patient successfully books a slot.
type Appointment struct {
	TimeSlotID  int
	PatientID   string
	PatientName string
	Age         int
	Phone       string
	Email       string
	Symptoms    string
}
