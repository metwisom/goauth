package user

import (
	"errors"
	"fmt"
	"goauth/internal/db"
	"goauth/internal/libs/steam"
	"goauth/internal/utils"
	"strings"
)

// Error definitions
var (
	ErrUserAlreadyExists = utils.NewMyError(409, "User already exists", "")
	ErrUserNotFound      = utils.NewMyError(404, "User not found", "")
	ErrInvalidInput      = utils.NewMyError(400, "Invalid input parameters", "")
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login,omitempty"`
	Username string `json:"username,omitempty"`
	SteamID  string `json:"steam_id,omitempty"`
}

// Create creates a new user with the specified login and password
// Returns the user ID or an error if the user cannot be created
func Create(login, password string) (id int, err error) {
	// Validate input
	login = strings.TrimSpace(login)
	if login == "" || password == "" {
		return 0, ErrInvalidInput
	}

	// Check if user already exists
	id, _, err = GetByLogin(login)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return 0, fmt.Errorf("failed to check existing user: %w", err)
	}
	if id != 0 {
		return 0, ErrUserAlreadyExists
	}

	// Hash password and create user
	hashPassword := utils.HashPassword(password)
	vars := []interface{}{login, hashPassword}
	err = db.DB.QueryRow("INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", vars, &id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

// CreateBySteam creates a new user from Steam authentication data
// Returns the user ID or an error if the user cannot be created
func CreateBySteam(steamUser steam.User) (id int, err error) {
	// Validate input
	if steamUser.SteamID == "" {
		return 0, ErrInvalidInput
	}

	// Check if user already exists
	existingID, err := GetBySteam(steamUser)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return 0, fmt.Errorf("failed to check existing Steam user: %w", err)
	}
	if existingID != 0 {
		return existingID, nil // Return existing user instead of error
	}

	// Create new user
	vars := []interface{}{steamUser.SteamID, steamUser.PersonaName}
	err = db.DB.QueryRow("INSERT INTO users (steam_id, username) VALUES ($1, $2) RETURNING id", vars, &id)
	if err != nil {
		return 0, fmt.Errorf("failed to create Steam user: %w", err)
	}

	return id, nil
}

// GetByLogin retrieves a user by login
// Returns user ID, password hash, and an error if the user was not found
func GetByLogin(login string) (id int, passwordHash string, err error) {
	login = strings.TrimSpace(login)
	if login == "" {
		return 0, "", ErrInvalidInput
	}

	vars := []interface{}{login}
	err = db.DB.QueryRow("SELECT id, password FROM users WHERE login = $1 LIMIT 1", vars, &id, &passwordHash)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return 0, "", ErrUserNotFound
		}
		return 0, "", fmt.Errorf("database error while getting user by login: %w", err)
	}
	return id, passwordHash, nil
}

// GetBySteam retrieves a user by Steam ID
// Returns user ID and an error if the user was not found
func GetBySteam(steamUser steam.User) (id int, err error) {
	// Validate input
	if steamUser.SteamID == "" {
		return 0, ErrInvalidInput
	}

	vars := []interface{}{steamUser.SteamID}
	err = db.DB.QueryRow("SELECT id FROM users WHERE steam_id = $1", vars, &id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("database error while getting user by Steam ID: %w", err)
	}

	return id, nil
}

// GetByID retrieves a user by ID
// Returns a User struct and an error if the user was not found
func GetByID(userID int) (*User, error) {
	if userID <= 0 {
		return nil, ErrInvalidInput
	}

	var user User
	vars := []interface{}{userID}
	err := db.DB.QueryRow("SELECT id, login, username, steam_id FROM users WHERE id = $1", vars,
		&user.ID, &user.Login, &user.Username, &user.SteamID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("database error while getting user by ID: %w", err)
	}

	return &user, nil
}
