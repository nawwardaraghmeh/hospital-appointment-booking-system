package repository

import (
	"abs/internal/entity"
	"database/sql"
)

// CreateAppointment 1. inserts the appointment record, and 2. marks the corresponding timeslot as booked
func CreateAppointment(db *sql.DB, a entity.Appointment) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 1. insert the appointment record
	_, err = tx.Exec(`
		INSERT INTO appointment(timeslot_id, patient_id, patient_name, age, phone, email, symptoms)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		a.TimeSlotID, a.PatientID, a.PatientName, a.Age, a.Phone, a.Email, a.Symptoms)

	if err != nil {
		tx.Rollback()
		return err
	}

	// 2. mark the corresponding timeslot as booked
	_, err = tx.Exec(`UPDATE timeslot SET is_booked=1 WHERE id=?`, a.TimeSlotID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
