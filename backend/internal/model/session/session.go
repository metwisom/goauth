package session

import (
	"errors"
	"goauth/internal/db"
	"goauth/internal/utils"
	"time"
)

func Create(userID int) (sessionID string, err error) {
	sessionID = utils.GenerateRandomString(30)
	expiresAt := time.Now().Add(24 * time.Hour)
	vars := []interface{}{sessionID, userID, expiresAt}
	err = db.DB.Exec("INSERT INTO sessions (session_id, user_id, expires_at) VALUES ($1, $2, $3)", vars)
	return
}

func GetUser(sessionID string) (userID int, err error) {
	var expiresAt time.Time
	vars := []interface{}{sessionID}
	err = db.DB.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_id = $1", vars, &userID, &expiresAt)
	if err != nil {
		return
	}
	if time.Now().After(expiresAt) {
		err = errors.New("session expired")
	}
	return
}
