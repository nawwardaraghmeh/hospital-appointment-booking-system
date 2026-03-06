package main

import (
	"abs/internal/config"
	"abs/internal/handler"
	"abs/internal/middleware"
	"abs/internal/repository"
	"abs/internal/service"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	// open and initialise the database
	db, err := repository.OpenDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer db.Close()

	if err := repository.InitDB(db); err != nil {
		log.Fatalf("schema error: %v", err)
	}
	if err := repository.SeedData(db); err != nil {
		log.Fatalf("seed error: %v", err)
	}

	// build repositories
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatalf("user repo: %v", err)
	}
	locationRepo, err := repository.NewLocationRepository(db)
	if err != nil {
		log.Fatalf("location repo: %v", err)
	}
	slotRepo, err := repository.NewTimeslotRepository(db)
	if err != nil {
		log.Fatalf("timeslot repo: %v", err)
	}
	apptRepo, err := repository.NewAppointmentRepository(db)
	if err != nil {
		log.Fatalf("appointment repo: %v", err)
	}

	// build services
	bookingSvc := service.NewBookingService(apptRepo)

	// build template renderer
	tmpl, err := handler.NewTemplateRenderer("templates")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	// build handlers
	authHandler := handler.NewAuthHandler(userRepo, locationRepo, cfg.AdminRegistrationCode, cfg.SessionDuration, tmpl)
	adminHandler := handler.NewAdminHandler(locationRepo, slotRepo, tmpl)
	patientHandler := handler.NewPatientHandler(userRepo, locationRepo, slotRepo, bookingSvc, tmpl)

	// register routes

	// public patient routes
	http.HandleFunc("/", authHandler.HomePage)
	http.HandleFunc("/login", authHandler.LoginPage)
	http.HandleFunc("/do-login", authHandler.LoginAction)
	http.HandleFunc("/register", authHandler.RegisterPage)
	http.HandleFunc("/do-register", authHandler.RegisterAction)
	http.HandleFunc("/logout", authHandler.Logout)

	// hidden admin registration. admins navigate here directly: /admin/register
	http.HandleFunc("/admin/register", authHandler.AdminRegisterPage)
	http.HandleFunc("/do-admin-register", authHandler.AdminRegisterAction)

	// protected — admin
	http.HandleFunc("/admin", middleware.RequireSession(adminHandler.Dashboard, "admin"))
	http.HandleFunc("/admin/add-slot", middleware.RequireSession(adminHandler.AddSlot, "admin"))

	// protected — patient
	http.HandleFunc("/patient/slots", middleware.RequireSession(patientHandler.SlotsPage, "patient"))
	http.HandleFunc("/book", middleware.RequireSession(patientHandler.BookPage, "patient"))

	// static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// start server
	addr := ":" + cfg.ServerPort
	fmt.Printf("Server running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
