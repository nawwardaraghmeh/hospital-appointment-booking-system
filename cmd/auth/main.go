package main

import (
	"abs/internal/bookingclient"
	"abs/internal/config"
	"abs/internal/handler"
	"abs/internal/middleware"
	"abs/internal/repository"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// open the shared SQLite database (Auth Service owns the users table)
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

	// build repositories that Auth Service needs (users only)
	userRepo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatalf("user repo error: %v", err)
	}

	// build the HTTP client that calls the Booking Service
	bookingClient := bookingclient.New(cfg.BookingURL)

	// parse and cache all HTML templates at startup
	tmpl, err := handler.NewTemplateRenderer("./templates")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	// inject dependencies into the single Auth handler
	authHandler := handler.NewAuthHandler(
		userRepo,
		bookingClient,
		cfg.AdminRegistrationCode,
		cfg.SessionDuration,
		tmpl,
	)

	mux := http.NewServeMux()

	// public routes
	mux.HandleFunc("/", authHandler.HomePage)
	mux.HandleFunc("/register", authHandler.RegisterPage)
	mux.HandleFunc("/do-register", authHandler.RegisterAction)
	mux.HandleFunc("/login", authHandler.LoginPage)
	mux.HandleFunc("/do-login", authHandler.LoginAction)
	mux.HandleFunc("/logout", authHandler.Logout)

	// hidden admin registration
	mux.HandleFunc("/admin/register", authHandler.AdminRegisterPage)
	mux.HandleFunc("/admin/do-register", authHandler.AdminRegisterAction)

	// protected routes: middleware checks the session cookie role
	mux.HandleFunc("/admin", middleware.RequireSession(authHandler.AdminDashboard, "admin"))
	mux.HandleFunc("/admin/add-slot", middleware.RequireSession(authHandler.AddSlot, "admin"))
	mux.HandleFunc("/patient/slots", middleware.RequireSession(authHandler.SlotsPage, "patient"))
	mux.HandleFunc("/book", middleware.RequireSession(authHandler.BookPage, "patient"))

	// static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	addr := ":" + cfg.AuthPort
	fmt.Printf("Auth Service running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
