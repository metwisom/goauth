package token

import (
	"goauth/internal/db"
	"goauth/internal/utils"
	"time"
)

func Create(userID int) (accessToken string, err error) {
	accessToken = utils.GenerateRandomString(32)
	expiresAt := time.Now().Add(1 * time.Hour)
	vars := []interface{}{accessToken, userID, expiresAt}
	err = db.DB.Exec("INSERT INTO access_token (access_token, user_id, expires_at) VALUES ($1, $2, $3)", vars)
	return
}
