package handler

import (
	"goauth/internal/config"
	"goauth/internal/httpServer/responseError"
	"goauth/internal/model/client"
	"goauth/internal/model/ot_code"
	"goauth/internal/model/session"
	"html"
	"log"
	"net/url"
	"strconv"

	"github.com/valyala/fasthttp"
)

type authorizeRequest struct {
	clientID     int
	redirectURI  string
	responseType string
	sessionID    string
}

// parseAndValidateRequest парсит и валидирует параметры запроса
func parseAndValidateRequest(ctx *fasthttp.RequestCtx) (*authorizeRequest, int, interface{}) {
	req := &authorizeRequest{
		redirectURI:  string(ctx.QueryArgs().Peek("redirect_uri")),
		responseType: string(ctx.QueryArgs().Peek("response_type")),
		sessionID:    string(ctx.Request.Header.Cookie(config.Config.SessionCookieName)),
	}

	// Проверка наличия обязательных параметров
	clientIDStr := string(ctx.QueryArgs().Peek("client_id"))
	if clientIDStr == "" || req.redirectURI == "" || req.responseType == "" {
		return nil, fasthttp.StatusBadRequest, responseError.BadRequest(responseError.ErrMissingParameters)
	}

	// Проверка сессии
	if req.sessionID == "" {
		ctx.Redirect(string(ctx.Request.URI().Host())+"/login?redirect_uri="+html.EscapeString(req.redirectURI)+"&response_type="+req.responseType+"&client_id="+clientIDStr, fasthttp.StatusFound)
		return nil, 302, responseError.BadRequest(responseError.ErrMissingCookie)
	}

	// Валидация client_id
	var err error
	req.clientID, err = strconv.Atoi(clientIDStr)
	if err != nil {
		return nil, fasthttp.StatusBadRequest, responseError.BadRequest(responseError.ErrClientIdError)
	}

	// Валидация redirect_uri
	if _, err := url.Parse(req.redirectURI); err != nil {
		return nil, fasthttp.StatusBadRequest, responseError.BadRequest(responseError.ErrInvalidRedirectUri)
	}

	// Валидация response_type
	if req.responseType != "code" {
		return nil, fasthttp.StatusBadRequest, responseError.BadRequest("unsupported response_type, only 'code' is allowed")
	}

	return req, 0, nil
}

// validateUserAndClient проверяет пользователя и клиента
func validateUserAndClient(req *authorizeRequest) (int, int, interface{}) {
	// Получение пользователя из сессии
	userID, err := session.GetUser(req.sessionID)
	if err != nil {
		log.Printf("Invalid session: %v", err)
		return 0, fasthttp.StatusBadRequest, responseError.BadRequest(responseError.ErrInvalidSession)
	}

	// Проверка клиента
	_, err = client.Get(req.clientID)
	if err != nil {
		log.Printf("Client not found: %v", err)
		return 0, fasthttp.StatusBadRequest, responseError.BadRequest("Client not found")
	}

	return userID, 0, nil
}

// createAuthorizationCode создает код авторизации и формирует URL для редиректа
func createAuthorizationCode(userID, clientID int, redirectURI string) (string, error) {
	code, err := ot_code.Create(userID, clientID)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(redirectURI)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("code", code.Code)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// Authorize обрабатывает запрос авторизации OAuth
func Authorize(ctx *fasthttp.RequestCtx) (int, interface{}) {
	// Проверка метода запроса
	if !ctx.IsGet() {
		return fasthttp.StatusMethodNotAllowed, responseError.MethodNotAllowed(responseError.ErrMethodNotAllowed)
	}

	// Парсинг и валидация запроса
	req, errCode, errResp := parseAndValidateRequest(ctx)
	if req == nil {
		return errCode, errResp
	}

	// Валидация пользователя и клиента
	userID, errCode, errResp := validateUserAndClient(req)
	if userID == 0 {
		return errCode, errResp
	}

	// Создание кода авторизации и формирование URL для редиректа
	redirectURL, err := createAuthorizationCode(userID, req.clientID, req.redirectURI)
	if err != nil {
		log.Printf("Failed to create authorization code: %v", err)
		return fasthttp.StatusInternalServerError, responseError.InternalServerError("Failed to create authorization code")
	}

	// Выполнение редиректа
	ctx.Redirect(redirectURL, fasthttp.StatusFound)
	return fasthttp.StatusFound, map[string]interface{}{}
}
