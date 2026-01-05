package repository

func SeedData() {
	db, _ := OpenDB()
	defer db.Close()

	cities := []string{"Rome", "Milan", "Naples", "Florence"}
	for _, name := range cities {
		db.Exec(`INSERT OR IGNORE INTO city(name) VALUES (?)`, name)
	}

	hospitals := []struct {
		name string
		city string
	}{
		{"San Giovanni", "Rome"}, {"Policlinico", "Rome"},
		{"Niguarda", "Milan"}, {"Fatebenefratelli", "Milan"},
		{"Ospedale Vecchio", "Naples"},
		{"Santa Maria Nuova", "Florence"}, {"Careggi", "Florence"},
	}

	for _, h := range hospitals {
		db.Exec(`INSERT OR IGNORE INTO hospital(name, city_id) 
                SELECT ?, id FROM city WHERE name = ?`, h.name, h.city)
	}

	depts := []string{"Cardiology", "Neurology", "Orthopedics", "Dermatology", "Pediatrics", "Oncology", "Radiology"}

	rows, _ := db.Query("SELECT id FROM hospital")
	var hospIDs []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		hospIDs = append(hospIDs, id)
	}
	rows.Close()

	for _, hID := range hospIDs {
		for _, dName := range depts {
			db.Exec(`INSERT OR IGNORE INTO department(name, hospital_id) VALUES (?, ?)`, dName, hID)
		}
	}
}
