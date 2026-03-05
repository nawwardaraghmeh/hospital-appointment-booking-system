package repository

import (
	"abs/internal/entity"
	"database/sql"
	"fmt"
)

// TimeslotRepository defines database operations for timeslots
type TimeslotRepository interface {
	FindAvailableByDepartment(departmentID string) ([]entity.Timeslot, error)
	FindAllByHospital(hospitalID string) ([]entity.Timeslot, error)
	FindDepartmentIDBySlotID(slotID int) (string, error)
	Create(departmentID, doctor, room, startTime string, duration int) error
}

type timeslotRepository struct {
	db                  *sql.DB
	stmtAvailableByDept *sql.Stmt
	stmtAllByHospital   *sql.Stmt
	stmtDeptIDBySlot    *sql.Stmt
	stmtCreate          *sql.Stmt
}

// NewTimeslotRepository constructs a TimeslotRepository with all statements prepared
func NewTimeslotRepository(db *sql.DB) (TimeslotRepository, error) {
	stmtAvail, err := db.Prepare(
		`SELECT id, doctor, room, start_time, duration
		 FROM timeslot
		 WHERE department_id = ? AND is_booked = 0`)
	if err != nil {
		return nil, fmt.Errorf("prepare available timeslots: %w", err)
	}

	stmtAll, err := db.Prepare(
		`SELECT t.id, d.name, t.doctor, t.room, t.start_time, t.duration,
		        t.is_booked, IFNULL(a.patient_name, '')
		 FROM timeslot t
		 JOIN department d        ON t.department_id = d.id
		 LEFT JOIN appointment a  ON t.id            = a.timeslot_id
		 WHERE d.hospital_id = ?`)
	if err != nil {
		return nil, fmt.Errorf("prepare all timeslots by hospital: %w", err)
	}

	stmtDept, err := db.Prepare(
		`SELECT department_id FROM timeslot WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("prepare dept id by slot: %w", err)
	}

	stmtCreate, err := db.Prepare(
		`INSERT INTO timeslot (department_id, doctor, room, start_time, duration, is_booked)
		 VALUES (?, ?, ?, ?, ?, 0)`)
	if err != nil {
		return nil, fmt.Errorf("prepare create timeslot: %w", err)
	}

	return &timeslotRepository{
		db:                  db,
		stmtAvailableByDept: stmtAvail,
		stmtAllByHospital:   stmtAll,
		stmtDeptIDBySlot:    stmtDept,
		stmtCreate:          stmtCreate,
	}, nil
}

func (r *timeslotRepository) FindAvailableByDepartment(departmentID string) ([]entity.Timeslot, error) {
	rows, err := r.stmtAvailableByDept.Query(departmentID)
	if err != nil {
		return nil, fmt.Errorf("find available timeslots: %w", err)
	}
	defer rows.Close()

	var slots []entity.Timeslot
	for rows.Next() {
		var s entity.Timeslot
		if err := rows.Scan(&s.ID, &s.Doctor, &s.Room, &s.StartTime, &s.Duration); err != nil {
			return nil, fmt.Errorf("scan timeslot: %w", err)
		}
		slots = append(slots, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate available timeslots: %w", err)
	}
	return slots, nil
}

func (r *timeslotRepository) FindAllByHospital(hospitalID string) ([]entity.Timeslot, error) {
	rows, err := r.stmtAllByHospital.Query(hospitalID)
	if err != nil {
		return nil, fmt.Errorf("find timeslots by hospital: %w", err)
	}
	defer rows.Close()

	var slots []entity.Timeslot
	for rows.Next() {
		var s entity.Timeslot
		var booked int
		if err := rows.Scan(&s.ID, &s.Department, &s.Doctor, &s.Room,
			&s.StartTime, &s.Duration, &booked, &s.Patient); err != nil {
			return nil, fmt.Errorf("scan timeslot row: %w", err)
		}
		s.IsBooked = booked == 1
		slots = append(slots, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate hospital timeslots: %w", err)
	}
	return slots, nil
}

func (r *timeslotRepository) FindDepartmentIDBySlotID(slotID int) (string, error) {
	var deptID string
	err := r.stmtDeptIDBySlot.QueryRow(slotID).Scan(&deptID)
	if err != nil {
		return "", fmt.Errorf("find department id for slot %d: %w", slotID, err)
	}
	return deptID, nil
}

func (r *timeslotRepository) Create(departmentID, doctor, room, startTime string, duration int) error {
	_, err := r.stmtCreate.Exec(departmentID, doctor, room, startTime, duration)
	if err != nil {
		return fmt.Errorf("create timeslot: %w", err)
	}
	return nil
}
