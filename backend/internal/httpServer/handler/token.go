package handler

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"goauth/internal/httpServer/responseError"
	"goauth/internal/model/client"
	"goauth/internal/model/ot_code"
	"goauth/internal/model/token"
	"log"
	"strconv"
)

type tokenRequest struct {
	code         string
	clientID     int
	clientSecret string
	grantType    string
}

// parseTokenRequest парсит и валидирует параметры запроса
func parseTokenRequest(ctx *fasthttp.RequestCtx) (*tokenRequest, int, interface{}) {
	req := &tokenRequest{
		code:         string(ctx.PostArgs().Peek("code")),
		clientSecret: string(ctx.PostArgs().Peek("client_secret")),
		grantType:    string(ctx.PostArgs().Peek("grant_type")),
	}
	fmt.Println(req.code)
	// Проверка наличия всех необходимых параметров
	if req.code == "" {
		return nil, fasthttp.StatusBadRequest, responseError.BadRequest("Missing required parameters1")
	}
	// Проверка наличия всех необходимых параметров
	if req.clientSecret == "" {
		return nil, fasthttp.StatusBadRequest, responseError.BadRequest("Missing required parameters2")
	}
	// Проверка наличия всех необходимых параметров
	if req.grantType == "" {
		return nil, fasthttp.StatusBadRequest, responseError.BadRequest("Missing required parameters3")
	}

	// Парсинг и валидация client_id
	clientIDStr := string(ctx.PostArgs().Peek("client_id"))
	if clientIDStr == "" {
		return nil, fasthttp.StatusBadRequest, responseError.BadRequest("Missing client_id")
	}

	var err error
	req.clientID, err = strconv.Atoi(clientIDStr)
	if err != nil {
		return nil, fasthttp.StatusBadRequest, responseError.BadRequest("Invalid client_id format")
	}

	return req, 0, nil
}

// validateAuthorizationCode проверяет код авторизации и клиента
func validateAuthorizationCode(req *tokenRequest) (*ot_code.Code, *client.Client, int, interface{}) {
	// Получение кода авторизации
	code, err := ot_code.Get(req.code)
	if err != nil {
		log.Printf("Failed to get authorization code: %v", err)
		return nil, nil, fasthttp.StatusBadRequest, responseError.BadRequest("Invalid authorization code")
	}

	// Получение и проверка клиента
	foundClient, err := client.Get(code.ClientID)
	if err != nil || foundClient.Secret == "" {
		log.Printf("Failed to get client or invalid client secret: %v", err)
		return nil, nil, fasthttp.StatusBadRequest, responseError.BadRequest("Invalid client")
	}

	// Проверка параметров запроса
	if req.grantType != "authorization_code" ||
		req.clientID != code.ClientID ||
		req.clientSecret != foundClient.Secret {
		return nil, nil, fasthttp.StatusBadRequest, responseError.BadRequest("Invalid request parameters")
	}

	return &code, &foundClient, 0, nil
}

// Token обрабатывает запрос на получение токена доступа
func Token(ctx *fasthttp.RequestCtx) (int, interface{}) {
	// Парсинг запроса
	req, errCode, errResp := parseTokenRequest(ctx)
	if errCode != 0 {
		return errCode, errResp
	}

	// Валидация кода и клиента
	code, _, errCode, errResp := validateAuthorizationCode(req)
	if errCode != 0 {
		return errCode, errResp
	}

	// Создание токена доступа
	accessToken, err := token.Create(code.UserID)
	if err != nil {
		log.Printf("Failed to create access token: %v", err)
		return fasthttp.StatusInternalServerError, responseError.InternalServerError("Failed to create access token")
	}

	// Удаление использованного кода
	if err := code.Delete(); err != nil {
		log.Printf("Failed to delete used authorization code: %v", err)
		return fasthttp.StatusInternalServerError, responseError.InternalServerError("Failed to process request")
	}

	// Формирование успешного ответа
	return fasthttp.StatusOK, map[string]interface{}{
		"access_token": accessToken,
		"token_type":   "bearer",
	}
}
