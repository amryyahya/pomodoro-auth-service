package utils

import (
	"database/sql"
	_ "github.com/lib/pq" 
)

func Blacklist_Token(db *sql.DB, user_id string, token string) (string, error) {
	query := `INSERT INTO token_blacklist (user_id, token, expiry_time) VALUES ($1, $2, $3)`
	_, err := db.Exec(user_id, token)
	if err != nil {
		return err
	}
	return nil
}