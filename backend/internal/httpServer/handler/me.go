package handler

import (
	"goauth/internal/config"
	"goauth/internal/httpServer/responseError"
	"goauth/internal/model/session"
	"log"

	"github.com/valyala/fasthttp"
)

// isPreflightRequest проверяет, является ли запрос предварительным OPTIONS запросом
func isPreflightRequest(ctx *fasthttp.RequestCtx) bool {
	return string(ctx.Method()) == "OPTIONS"
}

// getSessionID извлекает идентификатор сессии из кук
func getSessionID(ctx *fasthttp.RequestCtx) string {
	return string(ctx.Request.Header.Cookie(config.Config.SessionCookieName))
}

// GetMe обрабатывает запрос для получения информации о текущем пользователе.
// Возвращает HTTP статус и тело ответа.
func GetMe(ctx *fasthttp.RequestCtx) (int, interface{}) {

	if isPreflightRequest(ctx) {
		return fasthttp.StatusOK, nil
	}

	if !ctx.IsGet() {
		return fasthttp.StatusMethodNotAllowed, responseError.MethodNotAllowed(responseError.ErrMethodNotAllowed)
	}

	sessionID := getSessionID(ctx)
	if sessionID == "" {
		return fasthttp.StatusUnauthorized, responseError.Unauthorized("No session found")
	}

	userId, err := session.GetUser(sessionID)
	if err != nil {
		log.Printf("Failed to get user from session: %v", err)
		return fasthttp.StatusUnauthorized, responseError.Unauthorized("Invalid session")
	}

	return fasthttp.StatusOK, map[string]interface{}{
		"user_id": userId,
		"status":  "success",
	}
}
