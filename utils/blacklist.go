package utils

import (
	"database/sql"
	"time"
	_ "github.com/lib/pq" 
)

func BlacklistToken(db *sql.DB, user_id string, token string, expirationUnix float64) error {
	expiration := time.Unix(int64(expirationUnix), 0)
	query := `INSERT INTO token_blacklist (user_id, token, expiry_time) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, user_id, token, expiration)
	if err != nil {
		return err
	}
	return nil
}

func IsTokenBlacklisted(db *sql.DB, token string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM token_blacklist WHERE token = $1)`
	var exists bool
	err := db.QueryRow(query, token).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
