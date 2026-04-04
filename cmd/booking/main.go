package main

import (
	"abs/internal/bookingapi"
	"abs/internal/config"
	"abs/internal/repository"
	"abs/internal/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// open the shared SQLite database
	db, err := repository.OpenDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer db.Close()

	if err := repository.InitDB(db); err != nil {
		log.Fatalf("InitDB error: %v", err)
	}
	if err := repository.SeedData(db); err != nil {
		log.Fatalf("SeedData error: %v", err)
	}

	// build all repositories the Booking Service needs
	locationRepo, err := repository.NewLocationRepository(db)
	if err != nil {
		log.Fatalf("location repo error: %v", err)
	}
	slotRepo, err := repository.NewTimeslotRepository(db)
	if err != nil {
		log.Fatalf("timeslot repo error: %v", err)
	}
	apptRepo, err := repository.NewAppointmentRepository(db)
	if err != nil {
		log.Fatalf("appointment repo error: %v", err)
	}
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatalf("user repo error: %v", err)
	}

	bookingSvc := service.NewBookingService(apptRepo)
	_ = userRepo

	// wire all routes onto the Booking Service API handler
	apiHandler := bookingapi.New(locationRepo, slotRepo, bookingSvc)

	mux := http.NewServeMux()
	apiHandler.RegisterRoutes(mux)

	addr := ":" + cfg.BookingPort
	fmt.Printf("Booking Service running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
