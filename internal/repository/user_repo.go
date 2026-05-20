package repository

import (
	"abs/internal/entity"
	"database/sql"
	"fmt"
)

// define all database operations related to users
type UserRepository interface {
	Create(username, passwordHash, role, fullName string, hospitalID interface{}) error
	FindByUsername(username string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
}

type userRepository struct {
	db             *sql.DB
	stmtCreate     *sql.Stmt
	stmtByUsername *sql.Stmt
	stmtByID       *sql.Stmt
}

// create a UserRepository and prepare all statements
func NewUserRepository(db *sql.DB) (UserRepository, error) {
	stmtCreate, err := db.Prepare(
		`INSERT INTO users (username, password_hash, role, full_name, hospital_id)
		 VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("prepare user create: %w", err)
	}

	stmtByUsername, err := db.Prepare(
		`SELECT id, password_hash, role, IFNULL(hospital_id, 0)
		 FROM users WHERE username = ?`)
	if err != nil {
		return nil, fmt.Errorf("prepare user by username: %w", err)
	}

	stmtByID, err := db.Prepare(
		`SELECT id, full_name, IFNULL(hospital_id, 0) FROM users WHERE id = ?`)
	if err != nil {
		return nil, fmt.Errorf("prepare user by id: %w", err)
	}

	return &userRepository{
		db:             db,
		stmtCreate:     stmtCreate,
		stmtByUsername: stmtByUsername,
		stmtByID:       stmtByID,
	}, nil
}

// add a new user
func (r *userRepository) Create(username, passwordHash, role, fullName string, hospitalID interface{}) error {
	_, err := r.stmtCreate.Exec(username, passwordHash, role, fullName, hospitalID)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// find user by username, for login validation
func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	u := &entity.User{}
	err := r.stmtByUsername.QueryRow(username).Scan(&u.ID, &u.PasswordHash, &u.Role, &u.HospitalID)
	if err != nil {
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return u, nil
}

// find user by id, for token validation
func (r *userRepository) FindByID(id string) (*entity.User, error) {
	u := &entity.User{}
	err := r.stmtByID.QueryRow(id).Scan(&u.ID, &u.FullName, &u.HospitalID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return u, nil
}
