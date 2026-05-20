package repository

import "database/sql"

// populate the database with sample cities, hospitals, and departments
func SeedData(db *sql.DB) error {
	cities := []string{"Rome", "Milan", "Naples", "Florence"}
	for _, name := range cities {
		if _, err := db.Exec(`INSERT OR IGNORE INTO city(name) VALUES (?)`, name); err != nil {
			return err
		}
	}

	hospitals := []struct{ name, city string }{
		{"San Giovanni", "Rome"}, {"Policlinico", "Rome"},
		{"Niguarda", "Milan"}, {"Fatebenefratelli", "Milan"},
		{"Ospedale Vecchio", "Naples"},
		{"Santa Maria Nuova", "Florence"}, {"Careggi", "Florence"},
	}
	for _, h := range hospitals {
		_, err := db.Exec(
			`INSERT OR IGNORE INTO hospital(name, city_id)
			 SELECT ?, id FROM city WHERE name = ?`, h.name, h.city,
		)
		if err != nil {
			return err
		}
	}

	depts := []string{"Cardiology", "Neurology", "Orthopedics", "Dermatology", "Pediatrics", "Oncology", "Radiology"}
	rows, err := db.Query("SELECT id FROM hospital")
	if err != nil {
		return err
	}
	var hospIDs []int
	for rows.Next() {
		var id int
		rows.Scan(&id)
		hospIDs = append(hospIDs, id)
	}
	rows.Close()

	for _, hID := range hospIDs {
		for _, dName := range depts {
			_, err := db.Exec(`INSERT OR IGNORE INTO department(name, hospital_id) VALUES (?, ?)`, dName, hID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
