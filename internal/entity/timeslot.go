package entity

// this entity represents an available or booked appointment created by admin
type Timeslot struct {
	ID         int
	Department string
	Doctor     string
	Room       string
	StartTime  string
	Duration   int
	IsBooked   bool
	Patient    string
}
