package handler

import (
	"errors"
	"goauth/internal/httpServer/responseError"
	"goauth/internal/model/user"

	"github.com/valyala/fasthttp"
)

// RegisterRequest represents the expected request body for registration
type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Register handles user registration requests
// @Summary Register a new user
// @Description Creates a new user account with the provided login and password
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 200 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} responseError.Error "Invalid request parameters"
// @Failure 409 {object} responseError.Error "User already exists"
// @Router /register [post]
func Register(ctx *fasthttp.RequestCtx) (int, interface{}) {
	var req RegisterRequest

	req.Login = string(ctx.PostArgs().Peek("login"))
	req.Password = string(ctx.PostArgs().Peek("password"))

	// Validate input
	if req.Login == "" || req.Password == "" {
		return fasthttp.StatusBadRequest, responseError.BadRequest("Login and password are required")
	}

	if len(req.Login) < 3 || len(req.Login) > 32 {
		return fasthttp.StatusBadRequest, responseError.BadRequest("Login must be between 3 and 32 characters")
	}

	if len(req.Password) < 6 {
		return fasthttp.StatusBadRequest, responseError.BadRequest("Password must be at least 6 characters long")
	}

	// Create user
	newId, err := user.Create(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			return fasthttp.StatusConflict, responseError.BadRequest("User already exists")
		}
		return fasthttp.StatusInternalServerError, responseError.InternalServerError("Failed to create user")
	}

	return fasthttp.StatusOK, map[string]interface{}{
		"id":      newId,
		"message": "User created successfully",
	}
}
