# Hospital Appointment Booking Service (HABS)

A distributed web application built in Go for managing hospital timeslots and patient appointments. The system follows a **microservices architecture** with two independent services communicating over HTTP REST.

## Architecture

- **Auth Service** (port 8080): Handles user-facing concerns (HTML rendering, login, registration, sessions)
- **Booking Service** (port 8081): Owns hospital/appointment data, exposes JSON REST API
- **Database**: SQLite (shared between both services)

## Prerequisites

- Go 1.21 or higher

## Installation & Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd <project-directory>
```bash

2. Create the environment file

Create a .env file in the project root with the following content:

```bash
AUTH_PORT=8080
BOOKING_PORT=8081
DB_PATH=./data/abs.db
ADMIN_REGISTRATION_CODE=HOSPITAL2026
SESSION_DURATION=3600
BOOKING_SERVICE_URL=http://localhost:8081
```bash

3. Create data directory
```bash
mkdir -p data
```bash

4. Install dependencies
```bash
go mod tidy
```bash

5. Start both services

Terminal 1 - Booking Service (start first):
```bash
go run ./cmd/booking/main.go
```bash

Terminal 2 - Auth Service:
```bash
go run ./cmd/auth/main.go
```bash

On first run, the database file is created automatically, tables are initialised, and seed data (cities, hospitals, departments) is inserted.

6. Access the system
Open your browser and navigate to: http://localhost:8080

Reset Database
```bash
To reset the database at any time:
rm data/abs.db
```bash

Access Points
Home:	http://localhost:8080/	Public landing page
Patient: Registration	http://localhost:8080/register	Creates patient account
Patient/Admin Login:	http://localhost:8080/login	Role detected automatically
Admin Registration:	http://localhost:8080/admin/register	Hidden route - not linked in UI
Admin Dashboard:	http://localhost:8080/admin	After admin login
Patient Slots:	http://localhost:8080/patient/slots	After patient login

Admin registration requires the secret code defined in .env (ADMIN_REGISTRATION_CODE). The admin registration URL is intentionally not linked from the public UI, administrators navigate to it directly.