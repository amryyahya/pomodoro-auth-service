package models

import (
	"database/sql"
	_ "github.com/lib/pq" 
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CheckIfEmailExist(db *sql.DB, email string) (bool, error) {
	query := `SELECT 1 FROM users WHERE email = $1 LIMIT 1`
	var exists bool
	err := db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func InsertUser(db *sql.DB, email string, hashedPassword, salt []byte) error {
	query := `INSERT INTO users (email, hashed_password, salt) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, email, hashedPassword, salt)
	if err != nil {
		return err
	}
	return nil
}

func GetUserCredByEmail(db *sql.DB, email string) ([]byte, []byte, error) {
	// get user hashed_password and salt
	query := `SELECT hashed_password, salt FROM users WHERE email = $1`
	row := db.QueryRow(query, email)
	var hashedPassword, salt []byte
	if err := row.Scan(&hashedPassword, &salt); err != nil {
		return nil, nil, err
	}
	return hashedPassword, salt, nil
}

func GetUserIDByEmail(db *sql.DB, email string) (string, error) {
	query := `SELECT id FROM users WHERE email = $1`
	row := db.QueryRow(query, email)
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}
