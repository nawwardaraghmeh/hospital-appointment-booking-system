Appointment Booking System (ABS)
A distributed web application built in Go for managing hospital timeslots and patient appointments. 

Features
Role-Based Access Control: Separate flows for Admins and Patients using secure middleware.
Dynamic Filtering: Hierarchical selection (City → Hospital → Department) powered by JavaScript and Go backend.
Session Persistence: Custom session management using http.Cookie to track user state.
Concurrency Safety: Uses sync.Mutex in the booking service to prevent double-booking of slots.

Tech Stack
Language: Go (Golang)
Database: SQLite (via modernc.org/sqlite for a CGO-free implementation)
Frontend: HTML5, CSS3 (Admin/Patient styled consistently), and Vanilla JavaScript.

Database Schema
The system utilizes a relational schema to manage the hierarchy of healthcare facilities:
City: Primary location.
Hospital: Linked to a City.
Department: Linked to a Hospital (Standardized across all facilities).
Timeslot: Associated with a Department; tracks doctor and room availability.
Appointment: Links a Patient to a specific Timeslot.

Getting Started
Prerequisites
Go 1.21 or higher.

Installation & Run
Clone the repository.
Initialize the database and seed data:
go run cmd/server/main.go
Open http://localhost:8080 in your browser.
