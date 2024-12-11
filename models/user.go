package models

import (
	"database/sql"
)

type User struct {
    Email           string  `json:"email"`
    HashedPassword  string  `json:"hashedPassword"`
}

func InsertUser(db *sql.DB, email string, hashedPassword, salt []byte) error {
	query := `INSERT INTO users (email, hashed_password, salt) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, email, hashedPassword, salt)
	if err != nil {
		return err
	}
	return nil
}
