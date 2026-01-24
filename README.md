Appointment Booking System (ABS)
A distributed web application built in Go for managing hospital timeslots and patient appointments. 

Prerequisites
Go 1.21 or higher.

Installation & Run
1. Clone the repository.
2. Prepare the Environment, ensure you have a data folder in the root directory, otherwise run the following command
   mkdir data
3. Initialize & Start: Run the server.
   The first time you run this, the InitDB() and SeedData() functions called in main.go will automatically build your tables and populate the cities/hospitals.
   go run cmd/server/main.go
4. Access the systnen by openning the following link in your browser:
   http://localhost:8080
