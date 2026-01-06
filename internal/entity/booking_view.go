package entity

// this entity does not represent a single table, but the result of an SQL JOIN
// it's used to show patients data that concerns them
type BookingView struct {
	Doctor         string
	StartTime      string
	Room           string
	Symptoms       string
	PatientID      string
	DepartmentName string
}
