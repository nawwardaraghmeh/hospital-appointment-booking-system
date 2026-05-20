package repository

import (
	"abs/internal/entity"
	"database/sql"
	"fmt"
)

// handle read operations for cities, hospitals, and departments
type LocationRepository interface {
	AllCities() ([]entity.City, error)
	AllHospitals() ([]entity.Hospital, error)
	AllDepartments() ([]entity.Department, error)
	DepartmentsByHospital(hospitalID string) ([]entity.Department, error)
	FindDepartmentWithHospital(departmentID string) (deptName, hospitalName string, err error)
	FindHospitalNameByUserID(userID string) (adminName, hospitalName string, err error)
}

type locationRepository struct {
	db                     *sql.DB
	stmtAllCities          *sql.Stmt
	stmtAllHospitals       *sql.Stmt
	stmtAllDepartments     *sql.Stmt
	stmtDeptsByHospital    *sql.Stmt
	stmtDeptWithHospital   *sql.Stmt
	stmtHospitalNameByUser *sql.Stmt
}

// create a LocationRepository with all statements prepared
func NewLocationRepository(db *sql.DB) (LocationRepository, error) {
	stmts := map[string]string{
		"cities":       `SELECT id, name FROM city ORDER BY name`,
		"hospitals":    `SELECT id, name, city_id FROM hospital ORDER BY name`,
		"departments":  `SELECT id, name, hospital_id FROM department ORDER BY name`,
		"deptsByHosp":  `SELECT id, name FROM department WHERE hospital_id = ? ORDER BY name`,
		"deptWithHosp": `SELECT d.name, h.name FROM department d JOIN hospital h ON d.hospital_id = h.id WHERE d.id = ?`,
		"hospByUser":   `SELECT u.full_name, h.name FROM users u JOIN hospital h ON u.hospital_id = h.id WHERE u.id = ?`,
	}

	prepared := make(map[string]*sql.Stmt, len(stmts))
	for key, query := range stmts {
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, fmt.Errorf("prepare location stmt %q: %w", key, err)
		}
		prepared[key] = stmt
	}

	return &locationRepository{
		db:                     db,
		stmtAllCities:          prepared["cities"],
		stmtAllHospitals:       prepared["hospitals"],
		stmtAllDepartments:     prepared["departments"],
		stmtDeptsByHospital:    prepared["deptsByHosp"],
		stmtDeptWithHospital:   prepared["deptWithHosp"],
		stmtHospitalNameByUser: prepared["hospByUser"],
	}, nil
}

// retrieve all cities
func (r *locationRepository) AllCities() ([]entity.City, error) {
	rows, err := r.stmtAllCities.Query()
	if err != nil {
		return nil, fmt.Errorf("all cities: %w", err)
	}
	defer rows.Close()

	var cities []entity.City
	for rows.Next() {
		var c entity.City
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, fmt.Errorf("scan city: %w", err)
		}
		cities = append(cities, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate cities: %w", err)
	}
	return cities, nil
}

// retrieve all hospitals
func (r *locationRepository) AllHospitals() ([]entity.Hospital, error) {
	rows, err := r.stmtAllHospitals.Query()
	if err != nil {
		return nil, fmt.Errorf("all hospitals: %w", err)
	}
	defer rows.Close()

	var hospitals []entity.Hospital
	for rows.Next() {
		var h entity.Hospital
		if err := rows.Scan(&h.ID, &h.Name, &h.CityID); err != nil {
			return nil, fmt.Errorf("scan hospital: %w", err)
		}
		hospitals = append(hospitals, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate hospitals: %w", err)
	}
	return hospitals, nil
}

// retrieve all departments
func (r *locationRepository) AllDepartments() ([]entity.Department, error) {
	rows, err := r.stmtAllDepartments.Query()
	if err != nil {
		return nil, fmt.Errorf("all departments: %w", err)
	}
	defer rows.Close()

	var depts []entity.Department
	for rows.Next() {
		var d entity.Department
		if err := rows.Scan(&d.ID, &d.Name, &d.HospitalID); err != nil {
			return nil, fmt.Errorf("scan department: %w", err)
		}
		depts = append(depts, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate departments: %w", err)
	}
	return depts, nil
}

// retrieve departments for a specific hospital
func (r *locationRepository) DepartmentsByHospital(hospitalID string) ([]entity.Department, error) {
	rows, err := r.stmtDeptsByHospital.Query(hospitalID)
	if err != nil {
		return nil, fmt.Errorf("departments by hospital: %w", err)
	}
	defer rows.Close()

	var depts []entity.Department
	for rows.Next() {
		var d entity.Department
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, fmt.Errorf("scan department row: %w", err)
		}
		depts = append(depts, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate hospital departments: %w", err)
	}
	return depts, nil
}

// retrieve department name and its hospital name by department id
func (r *locationRepository) FindDepartmentWithHospital(departmentID string) (string, string, error) {
	var deptName, hospName string
	err := r.stmtDeptWithHospital.QueryRow(departmentID).Scan(&deptName, &hospName)
	if err != nil {
		return "", "", fmt.Errorf("find department+hospital: %w", err)
	}
	return deptName, hospName, nil
}

// retrieve hospital name and admin name by user id
func (r *locationRepository) FindHospitalNameByUserID(userID string) (string, string, error) {
	var adminName, hospitalName string
	err := r.stmtHospitalNameByUser.QueryRow(userID).Scan(&adminName, &hospitalName)
	if err != nil {
		return "", "", fmt.Errorf("find hospital name by user: %w", err)
	}
	return adminName, hospitalName, nil
}
