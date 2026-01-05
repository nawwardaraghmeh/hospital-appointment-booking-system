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
	repository.InitDB()
	repository.SeedData()

	http.HandleFunc("/patient/login", handler.PatientLoginPage)
	http.HandleFunc("/admin/login", handler.AdminLoginPage)

	http.HandleFunc("/admin/do-login", handler.AdminLoginAction)
	http.HandleFunc("/patient/do-login", handler.PatientLoginAction)

	http.HandleFunc("/admin/dashboard", middleware.RequireSession(handler.AdminPage, "admin"))
	http.HandleFunc("/patient/slots", middleware.RequireSession(handler.PatientSlotsPage, "patient"))

	http.HandleFunc("/admin/add-slot", middleware.RequireSession(handler.AddSlot, "admin"))
	http.HandleFunc("/book", middleware.RequireSession(handler.BookAppointment, "patient"))

	http.HandleFunc("/admin/login-submit", handler.AdminLoginPost)

	http.HandleFunc("/", handler.IndexPage)
	http.HandleFunc("/logout", handler.LogoutPage)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
