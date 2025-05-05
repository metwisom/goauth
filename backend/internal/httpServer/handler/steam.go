package handler

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"goauth/internal/config"
	"goauth/internal/httpServer/responseError"
	"goauth/internal/libs/steam"
	"goauth/internal/model/session"
	"goauth/internal/model/user"
	"goauth/internal/utils"
	"html"
)

func Steam(ctx *fasthttp.RequestCtx) (int, interface{}) {
	// Извлекаем nonce из кук
	redirectURI, err := validateRedirectURI(string(ctx.QueryArgs().Peek("redirect_uri")))
	responseType := string(ctx.QueryArgs().Peek("response_type"))
	clientId := string(ctx.QueryArgs().Peek("client_id"))

	nonce := string(ctx.Request.Header.Cookie("steam_nonce"))
	if nonce == "" {
		return fasthttp.StatusUnauthorized, responseError.Unauthorized("Invalid session")
	}

	// Парсим query параметры из URL
	claimId := ""
	query := fasthttp.Args{}
	ctx.URI().QueryArgs().VisitAll(func(key, value []byte) {
		if string(key) == "openid.claimed_id" {
			claimId = string(value)
		}
		query.Add(string(key), string(value))
	})

	valid, err := steam.ValidateSteamResponse(&query)
	if err != nil || !valid {
		return fasthttp.StatusUnauthorized, responseError.Unauthorized("Invalid Steam response")
	}

	// Извлекаем Steam ID
	claimedID := claimId
	steamID := steam.ExtractSteamID(claimedID)
	if steamID == "" {
		return fasthttp.StatusBadRequest, responseError.BadRequest("Invalid Steam ID")
	}

	// Получаем данные пользователя из Steam API
	foundUser, err := steam.GetSteamUser(steamID)
	if err != nil {
		return fasthttp.StatusInternalServerError, responseError.InternalServerError("Failed to fetch Steam user")
	}

	userId, err := user.GetBySteam(foundUser)

	if userId == 0 {
		userId, err = user.CreateBySteam(foundUser)
		fmt.Println(userId)
		if err != nil {
			fmt.Println(err)
		}
	}

	sessionID, err := session.Create(userId)
	ctx.Response.Header.SetCookie(utils.CreateCookie(config.Config.SessionCookieName, sessionID, '/'))

	fmt.Println("redirectURI:", redirectURI)
	fmt.Println("responseType:", responseType)
	fmt.Println("clientId:", clientId)

	if redirectURI == "" {
		ctx.Redirect("/authorize", fasthttp.StatusFound)
		return fasthttp.StatusFound, nil
	} else if redirectURI == "/" {
		ctx.Redirect(redirectURI, fasthttp.StatusFound)
		return fasthttp.StatusFound, map[string]interface{}{}
	} else {
		ctx.Redirect(string(ctx.Request.URI().Host())+"/api/authorize?redirect_uri="+html.EscapeString(redirectURI)+"&response_type="+responseType+"&client_id="+clientId, fasthttp.StatusFound)
		return fasthttp.StatusFound, map[string]interface{}{}
	}

	//// Генерируем JWT токен
	//token, err := generateJWT(user.SteamID)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Failed to generate JWT: %v\n", err)
	//	ctx.Error("Failed to generate token", fasthttp.StatusInternalServerError)
	//	return
	//}
	//
	//// Устанавливаем куки с токеном
	//ctx.Response.Header.SetCookie(&fasthttp.Cookie{
	//	Name:     "token",
	//	Value:    token,
	//	Path:     "/",
	//	HttpOnly: true,
	//	Expires:  time.Now().Add(1 * time.Hour),
	//})
	//
	//// Устанавливаем заголовок для JSON ответа
	//ctx.SetContentType("application/json")
	//_, err = ctx.WriteString(`{"message": "Authentication successful", "steam_id": "` + steamID + `"}`)
	//if err != nil {
	//	ctx.Error("Error sending response", fasthttp.StatusInternalServerError)
	//	return
	//}
	//
	//// Устанавливаем статус успешного ответа
}
