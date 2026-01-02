package entity

type City struct {
	ID   int
	Name string
}

type Hospital struct {
	ID     int
	Name   string
	CityID int
}
