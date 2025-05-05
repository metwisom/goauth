package handler

import (
	"fmt"
	"goauth/internal/config"
	"goauth/internal/httpServer/responseError"
	"goauth/internal/model/session"
	"goauth/internal/model/user"
	"goauth/internal/utils"
	"html"
	"log"
	"net/url"

	"github.com/valyala/fasthttp"
)

type loginRequest struct {
	login        string
	password     string
	redirectURI  string
	responseType string
	clientId     string
}

// validateRedirectURI проверяет корректность URI для перенаправления
func validateRedirectURI(uri string) (string, error) {
	if uri == "" {
		return "/", nil
	}

	u, err := url.Parse(uri)
	if err != nil || u.Path == "" {
		return "", utils.NewMyError(400, "invalid redirect URI", "")
	}

	return uri, nil
}

// parseLoginRequest извлекает и валидирует параметры запроса
func parseLoginRequest(ctx *fasthttp.RequestCtx) (*loginRequest, error) {
	redirectURI, err := validateRedirectURI(string(ctx.QueryArgs().Peek("redirect_uri")))
	responseType := string(ctx.QueryArgs().Peek("response_type"))
	clientId := string(ctx.QueryArgs().Peek("client_id"))
	if err != nil {
		return nil, err
	}

	req := &loginRequest{
		login:        string(ctx.PostArgs().Peek("login")),
		password:     string(ctx.PostArgs().Peek("password")),
		redirectURI:  redirectURI,
		responseType: responseType,
		clientId:     clientId,
	}

	if req.login == "" || req.password == "" {
		return nil, utils.NewMyError(400, responseError.ErrMissingCredentials, "")
	}

	return req, nil
}

// authenticateUser проверяет учетные данные пользователя
func authenticateUser(login, password string) (int, error) {
	userID, passwordHash, err := user.GetByLogin(login)
	if err != nil {
		fmt.Println(err)
		return 0, utils.NewMyError(400, "Нет юзера", "")
	}

	if !utils.CheckPassword(passwordHash, password) {
		return 0, utils.NewMyError(400, "Не тот пароль", "")
	}

	return userID, nil
}

// createSessionAndSetCookie создает сессию и устанавливает куки
func createSessionAndSetCookie(ctx *fasthttp.RequestCtx, userID int) error {
	sessionID, err := session.Create(userID)
	if err != nil {
		return err
	}

	cookie := utils.CreateCookie(
		config.Config.SessionCookieName,
		sessionID,
		86400, // 24 hours in seconds
	)
	ctx.Response.Header.SetCookie(cookie)
	return nil
}

// Login обрабатывает POST-запросы для входа пользователя
func Login(ctx *fasthttp.RequestCtx) (int, interface{}) {
	// Проверка метода запроса
	if !ctx.IsPost() {
		return fasthttp.StatusMethodNotAllowed, responseError.BadRequest(responseError.ErrMethodNotAllowed)
	}

	// Парсинг и валидация запроса
	req, err := parseLoginRequest(ctx)
	if err != nil {
		return fasthttp.StatusBadRequest, responseError.BadRequest(err.Error())
	}

	// Аутентификация пользователя
	userID, err := authenticateUser(req.login, req.password)
	if err != nil {
		log.Printf("Authentication failed for user %s from IP %s: %v", req.login, ctx.RemoteIP(), err)
		return fasthttp.StatusUnauthorized, responseError.Unauthorized(responseError.ErrInvalidCredentials)
	}

	// Создание сессии и установка куки
	if err := createSessionAndSetCookie(ctx, userID); err != nil {
		log.Printf("Session creation failed for user %d from IP %s: %v", userID, ctx.RemoteIP(), err)
		return fasthttp.StatusInternalServerError, responseError.InternalServerError(responseError.ErrSessionCreationFail)
	}

	if req.redirectURI == "/" {
		ctx.Redirect(req.redirectURI, fasthttp.StatusFound)
		return fasthttp.StatusFound, map[string]interface{}{}
	} else {
		ctx.Redirect(string(ctx.Request.URI().Host())+"/api/authorize?redirect_uri="+html.EscapeString(req.redirectURI)+"&response_type="+req.responseType+"&client_id="+req.clientId, fasthttp.StatusFound)
		return fasthttp.StatusFound, map[string]interface{}{}
	}

	// Перенаправление пользователя
}
