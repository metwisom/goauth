package handler

import (
	"github.com/valyala/fasthttp"
	"goauth/internal/config"
)

// Logout Обработчик для /logout
func Logout(ctx *fasthttp.RequestCtx) (int, interface{}) {

	redirectURI := string(ctx.QueryArgs().Peek("redirect_uri"))

	deleteCookie := fasthttp.AcquireCookie()
	deleteCookie.SetKey(config.Config.SessionCookieName)
	deleteCookie.SetMaxAge(-1)
	deleteCookie.SetPath("/")
	ctx.Response.Header.SetCookie(deleteCookie)

	ctx.Redirect(redirectURI, fasthttp.StatusFound)
	return fasthttp.StatusFound, nil

}
