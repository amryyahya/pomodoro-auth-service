package models

import (
	"database/sql"
)

type User struct {
    Email           string  `json:"email"`
    Password	 	string  `json:"password"`
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
