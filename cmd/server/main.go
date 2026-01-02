package main

import (
	"fmt"
	"log"
	"net/http"

	"abs/internal/handler"
	"abs/internal/repository"
)

func main() {
	repository.InitDB()
	repository.SeedData()

	fmt.Println("Current Timeslots in DB:")
	repository.PrintTimeslots()

	http.HandleFunc("/", handler.IndexPage)
	http.HandleFunc("/logout", handler.LogoutPage)

	// Patient
	http.HandleFunc("/patient/login", handler.PatientLoginPage)
	http.HandleFunc("/patient/slots", handler.PatientSlotsPage)
	http.HandleFunc("/book", handler.BookAppointment)

	// Admin
	http.HandleFunc("/admin/login", handler.AdminLoginPage)
	http.HandleFunc("/admin/dashboard", handler.AdminPage)
	http.HandleFunc("/admin/add-slot", handler.AddSlot)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
