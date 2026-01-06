package service

import (
	"abs/internal/entity"
	"abs/internal/repository"
	"database/sql"
	"sync"
)

var mutex sync.Mutex

// Reserve prevents "Race Conditions" where two patients might try to book the same timeslot at the exact same time
func Reserve(db *sql.DB, a entity.Appointment) error {
	mutex.Lock()
	defer mutex.Unlock()
	return repository.CreateAppointment(db, a)
}
