package repository

import (
	"abs/internal/entity"
	"database/sql"
)

func CreateAppointment(db *sql.DB, a entity.Appointment) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO appointment(timeslot_id, patient_id, patient_name, age, phone, email, symptoms)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		a.TimeSlotID, a.PatientID, a.PatientName, a.Age, a.Phone, a.Email, a.Symptoms)

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`UPDATE timeslot SET is_booked=1 WHERE id=?`, a.TimeSlotID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
