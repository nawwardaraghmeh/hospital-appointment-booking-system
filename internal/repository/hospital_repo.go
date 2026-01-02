package repository

import "database/sql"

func GetCities(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM city")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var c string
		rows.Scan(&c)
		cities = append(cities, c)
	}
	return cities, nil
}
