package repository

import (
	"abs/internal/entity"
	"database/sql"
)

func GetAvailableSlots(db *sql.DB, deptID int) ([]entity.TimeSlot, error) {
	rows, err := db.Query(`
		SELECT id, doctor, room, start_time, duration
		FROM timeslot
		WHERE department_id=? AND is_booked=0`, deptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []entity.TimeSlot
	for rows.Next() {
		var s entity.TimeSlot
		rows.Scan(&s.ID, &s.Doctor, &s.Room, &s.StartTime, &s.Duration)
		slots = append(slots, s)
	}
	return slots, nil
}
