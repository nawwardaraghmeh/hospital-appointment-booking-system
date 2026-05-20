package service

import (
	"abs/internal/entity"
	"abs/internal/repository"
	"sync"
)

// define database operations for booking
type BookingService interface {
	Reserve(a entity.Appointment) error
	MyAppointments(patientID string) ([]entity.BookingView, error)
}

type bookingService struct {
	apptRepo repository.AppointmentRepository
	mu       sync.Mutex
}

// create a BookingService with the given AppointmentRepository
func NewBookingService(apptRepo repository.AppointmentRepository) BookingService {
	return &bookingService{apptRepo: apptRepo}
}

// get a mutex lock before writing to prevent race conditions
func (s *bookingService) Reserve(a entity.Appointment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.apptRepo.Create(a)
}

// delegate to the repository to fetch a patient's bookings
func (s *bookingService) MyAppointments(patientID string) ([]entity.BookingView, error) {
	return s.apptRepo.FindByPatientID(patientID)
}
