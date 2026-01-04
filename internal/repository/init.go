package repository

func SeedData() {
	db, _ := OpenDB()
	defer db.Close()

	db.Exec(`INSERT OR IGNORE INTO city(id,name) VALUES (1,'Rome'),(2,'Milan'),(3,'Naples')`)

	db.Exec(`INSERT OR IGNORE INTO hospital(id,name,city_id) VALUES
		(1,'San Giovanni',1),
		(2,'Policlinico',1),
		(3,'Niguarda',2),
		(4,'Fatebenefratelli',2),
		(5,'Ospedale Vecchio',3)
	`)

	db.Exec(`INSERT OR IGNORE INTO department(id,name,hospital_id) VALUES
		(1,'Cardiology',1),
		(2,'Neurology',1),
		(3,'Orthopedics',2),
		(4,'Dermatology',2),
		(5,'Pediatrics',3),
		(6,'Oncology',4),
		(7,'Radiology',5)
	`)

	db.Exec(`INSERT OR IGNORE INTO admin(emp_id,name,city_id,hospital_id) VALUES
		('EMP001','Admin One',1,1)
	`)
}
