package repository

import (
	"abs/internal/entity"
	"database/sql"
	"errors"
	"fmt"
)

// AppointmentRepository defines database operations for appointments
type AppointmentRepository interface {
	Create(a entity.Appointment) error
	FindByPatientID(patientID string) ([]entity.BookingView, error)
}

type appointmentRepository struct {
	db                *sql.DB
	stmtFindByPatient *sql.Stmt
}

// NewAppointmentRepository creates an AppointmentRepository with prepared statements
func NewAppointmentRepository(db *sql.DB) (AppointmentRepository, error) {
	stmtFind, err := db.Prepare(
		`SELECT t.doctor, t.start_time, t.room, a.symptoms, d.name
		 FROM appointment a
		 JOIN timeslot   t ON a.timeslot_id   = t.id
		 JOIN department d ON t.department_id = d.id
		 WHERE a.patient_id = ?`)
	if err != nil {
		return nil, fmt.Errorf("prepare find appointments by patient: %w", err)
	}

	return &appointmentRepository{
		db:                db,
		stmtFindByPatient: stmtFind,
	}, nil
}

// Create atomically claims the timeslot (only if not already booked) and inserts the appointment
func (r *appointmentRepository) Create(a entity.Appointment) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.Exec(`UPDATE timeslot SET is_booked=1 WHERE id=? AND is_booked=0`, a.TimeSlotID)
	if err != nil {
		return fmt.Errorf("claim slot: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("claim slot rows affected: %w", err)
	}
	if rows == 0 {
		return errors.New("slot already booked")
	}

	_, err = tx.Exec(
		`INSERT INTO appointment(timeslot_id, patient_id, patient_name, age, phone, email, symptoms)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		a.TimeSlotID, a.PatientID, a.PatientName, a.Age, a.Phone, a.Email, a.Symptoms,
	)
	if err != nil {
		return fmt.Errorf("insert appointment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// FindByPatientID returns all appointments booked by a given patient
func (r *appointmentRepository) FindByPatientID(patientID string) ([]entity.BookingView, error) {
	rows, err := r.stmtFindByPatient.Query(patientID)
	if err != nil {
		return nil, fmt.Errorf("find appointments by patient: %w", err)
	}
	defer rows.Close()

	var results []entity.BookingView
	for rows.Next() {
		var bv entity.BookingView
		if err := rows.Scan(&bv.Doctor, &bv.StartTime, &bv.Room, &bv.Symptoms, &bv.DepartmentName); err != nil {
			return nil, fmt.Errorf("scan booking view: %w", err)
		}
		results = append(results, bv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate appointments: %w", err)
	}
	return results, nil
}
