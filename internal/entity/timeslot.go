package entity

type TimeSlot struct {
	ID           int
	DepartmentID int
	Doctor       string
	Room         string
	StartTime    string
	Duration     int
	IsBooked     bool
}
