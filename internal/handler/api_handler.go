package handler

import (
	"abs/internal/repository"
	"encoding/json"
	"net/http"
)

func CitiesAPI(w http.ResponseWriter, r *http.Request) {
	db, _ := repository.OpenDB()
	defer db.Close()

	rows, _ := db.Query("SELECT id, name FROM city")
	var data []map[string]interface{}

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		data = append(data, map[string]interface{}{
			"id": id, "name": name,
		})
	}
	json.NewEncoder(w).Encode(data)
}

func HospitalsAPI(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	db, _ := repository.OpenDB()
	defer db.Close()

	rows, _ := db.Query("SELECT id, name FROM hospital WHERE city_id=?", city)
	var data []map[string]interface{}

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		data = append(data, map[string]interface{}{
			"id": id, "name": name,
		})
	}
	json.NewEncoder(w).Encode(data)
}

func DepartmentsAPI(w http.ResponseWriter, r *http.Request) {
	h := r.URL.Query().Get("hospital")
	db, _ := repository.OpenDB()
	defer db.Close()

	rows, _ := db.Query("SELECT id, name FROM department WHERE hospital_id=?", h)
	var data []map[string]interface{}

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		data = append(data, map[string]interface{}{
			"id": id, "name": name,
		})
	}
	json.NewEncoder(w).Encode(data)
}
