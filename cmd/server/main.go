package main

import (
	"abs/internal/handler"
	"abs/internal/middleware"
	"abs/internal/repository"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// database setup
	repository.InitDB()
	repository.SeedData()

	// public routes
	http.HandleFunc("/", handler.HomePage)
	http.HandleFunc("/register", handler.RegisterPage)
	http.HandleFunc("/do-register", handler.RegisterAction)

	// specific login views for admin/patient
	http.HandleFunc("/patient/login", handler.PatientLoginPage)
	http.HandleFunc("/admin/login", handler.AdminLoginPage)

	// universal login
	http.HandleFunc("/do-login", handler.LoginAction)

	// protected Admin Routes
	http.HandleFunc("/admin", middleware.RequireSession(handler.AdminPage, "admin"))
	http.HandleFunc("/admin/add-slot", middleware.RequireSession(handler.AddSlot, "admin"))

	// protected Patient Routes
	http.HandleFunc("/patient/slots", middleware.RequireSession(handler.PatientSlotsPage, "patient"))
	http.HandleFunc("/book", middleware.RequireSession(handler.BookAppointment, "patient"))

	// session Termination
	http.HandleFunc("/logout", handler.Logout)

	// static Files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
