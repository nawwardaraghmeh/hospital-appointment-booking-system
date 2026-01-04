package entity

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
