package utils

import (
	"github.com/valyala/fasthttp"
	"strings"
)

func CreateCookie(key string, value string, expire int) *fasthttp.Cookie {
	if strings.Compare(key, "") == 0 {
		key = "ErrorToken"
	}
	authCookie := fasthttp.AcquireCookie()
	authCookie.SetKey(key)
	authCookie.SetPath("/")
	authCookie.SetValue(value)
	authCookie.SetMaxAge(expire)
	authCookie.SetHTTPOnly(true)
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	return authCookie
}
