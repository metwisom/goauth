package client

import (
	"goauth/internal/db"
	"goauth/internal/utils"
)

type Client struct {
	ClientID int    `db:"client_id"`
	UserID   int    `db:"user_id"`
	Secret   string `db:"secret"`
}

func Create(userID int) (client Client, err error) {
	client = Client{
		ClientID: 0,
		UserID:   userID,
		Secret:   utils.GenerateRandomString(16),
	}
	vars := []interface{}{client.UserID, client.Secret}
	err = db.DB.QueryRow("INSERT INTO client (user_id, secret) VALUES ($1, $2) RETURNING client_id", vars, &client.ClientID)
	return
}

func Get(clientId int) (client Client, err error) {
	vars := []interface{}{clientId}
	err = db.DB.QueryRow("SELECT client_id, user_id, secret FROM client WHERE client_id = $1", vars, &client.ClientID, &client.UserID, &client.Secret)
	return
}
