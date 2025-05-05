package ot_code

import (
	"errors"
	"goauth/internal/db"
	"goauth/internal/utils"
	"time"
)

type Code struct {
	ClientID int       `db:"client_id"`
	Code     string    `db:"code"`
	UserID   int       `db:"user_id"`
	ExpireAt time.Time `db:"expire_at"`
}

func (c *Code) Delete() (err error) {
	vars := []interface{}{c.Code}
	err = db.DB.Exec("DELETE from codes WHERE code = $1", vars)
	return
}

func Create(userID int, clientID int) (code Code, err error) {
	code = Code{
		Code:     utils.GenerateRandomString(18),
		UserID:   userID,
		ClientID: clientID,
		ExpireAt: time.Now().Add(1 * time.Hour),
	}
	vars := []interface{}{code.Code, code.UserID, code.ClientID, code.ExpireAt}
	err = db.DB.Exec("INSERT INTO codes (code, user_id, client_id, expires_at) VALUES ($1, $2, $3, $4)", vars)
	return
}

func Get(codeStr string) (code Code, err error) {
	vars := []interface{}{codeStr}
	err = db.DB.QueryRow("SELECT user_id, client_id, expires_at, code FROM codes WHERE code = $1", vars, &code.UserID, &code.ClientID, &code.ExpireAt, &code.Code)
	if err != nil {
		return
	}
	if time.Now().After(code.ExpireAt) {
		err = errors.New("code expire")
	}
	return
}
