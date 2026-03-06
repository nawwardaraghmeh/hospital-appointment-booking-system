**Appointment Booking System (ABS)**
A distributed web application built in Go for managing hospital timeslots and patient appointments. 

**Prerequisites**
Go 1.21 or higher.

**Installation & Run**
1. Clone the repository.
<<<<<<< HEAD
2. Create the environment file in the project root:
   `cp .env.example .env`
   Or create .env manually with the following content:
   `SERVER_PORT=8080
   DB_PATH=./data/abs.db
   ADMIN_REGISTRATION_CODE=HOSPITAL2026
   SESSION_DURATION=3600`
3. Initialize & Start: Run the server. On first run, the database tables and seed data are created automatically.
   `go run cmd/server/main.go`
4. Access the system at:
=======
2. Prepare the Environment, ensure you have a data folder in the root directory. If it's missing, run the following command
   `mkdir data`
3. Fetch dependencies: run the following command to install the required libraries
   `go mod tidy`
4. Initialize & Start: Run the server.
   The first time you run this, the InitDB() and SeedData() functions called in main.go will automatically build your tables and populate the cities/hospitals.
   `go run cmd/server/main.go`
5. Access the system by openning the following link in your browser:
>>>>>>> 486ea31236f190e16d3e93ae11b515f0b7da0b90
   http://localhost:8080

To reset the database at any time: `rm data/abs.db`

**Access**
Patient Register: /register
Patient Login: /login 
Admin Register: /admin/register
Admin Login: /login

Admin registration requires the secret code defined in .env.
After login, users are redirected automatically based on their role.
