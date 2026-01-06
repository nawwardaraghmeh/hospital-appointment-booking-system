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
	// setup schema
	repository.InitDB()
	// populate db
	repository.SeedData()

	// public routes: ones accessible by anyone
	http.HandleFunc("/", handler.IndexPage)
	http.HandleFunc("/patient/login", handler.PatientLoginPage)
	http.HandleFunc("/admin/login", handler.AdminLoginPage)

	// auth actions: handle POST requests
	http.HandleFunc("/admin/do-login", handler.AdminLoginAction)
	http.HandleFunc("/patient/do-login", handler.PatientLoginAction)

	// protected admin routes
	http.HandleFunc("/admin/dashboard", middleware.RequireSession(handler.AdminPage, "admin"))
	http.HandleFunc("/admin/add-slot", middleware.RequireSession(handler.AddSlot, "admin"))

	// protected patient routes
	http.HandleFunc("/patient/slots", middleware.RequireSession(handler.PatientSlotsPage, "patient"))
	http.HandleFunc("/book", middleware.RequireSession(handler.BookAppointment, "patient"))

	// session termination
	http.HandleFunc("/logout", handler.LogoutPage)

	// server setup
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
