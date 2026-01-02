package service

import (
	"abs/internal/entity"
	"abs/internal/repository"
	"database/sql"
	"sync"
)

var mutex sync.Mutex

func Reserve(db *sql.DB, a entity.Appointment) error {
	mutex.Lock()
	defer mutex.Unlock()
	return repository.CreateAppointment(db, a)
}
