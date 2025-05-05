package httpServer

import (
	"encoding/json"
	"fmt"
	"goauth/internal/httpServer/handler"
	"os"

	"github.com/valyala/fasthttp"
)

func CreateHTTPServer() error {
	isDevMode := os.Getenv("APP_ENV") == "dev"
	var proxyClient *fasthttp.HostClient

	// Если dev-режим, настраиваем прокси
	if isDevMode {
		proxyClient = &fasthttp.HostClient{
			Addr: "localhost:3000", // Адрес React dev server
		}
		fmt.Println("Запущен в dev-режиме, проксирование на localhost:3000")
	} else {
		fmt.Println("Запущен в prod-режиме, отдача файлов из ../frontend/build")
	}

	fs := &fasthttp.FS{
		Root:               "../frontend/build",    // Путь к папке относительно корня проекта
		IndexNames:         []string{"index.html"}, // Отдавать index.html по умолчанию
		GenerateIndexPages: false,                  // Не генерировать список файлов (обычно не нужно для React)
		Compress:           true,                   // Сжимать файлы для оптимизации
		AcceptByteRange:    true,                   // Поддержка частичной загрузки
		PathRewrite: func(ctx *fasthttp.RequestCtx) []byte {
			path := ctx.Path()
			return path
		},
	}

	// Создаем обработчик для статических файлов
	fsHandler := fs.NewRequestHandler()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		//aa := time.Now().UnixNano()
		//
		//defer func(aa int64) {
		//	fmt.Println("Close")
		//	fmt.Println(time.Now().UnixNano() - aa)
		//}(aa)

		switch string(ctx.Path()) {
		case "/api/login":

			//start := time.Now().UnixNano()

			ctx.SetContentType("application/json")
			code, data := handler.Login(ctx)
			ctx.SetStatusCode(code)
			err := json.NewEncoder(ctx).Encode(data)
			if err != nil {
				return
			}

			//fmt.Println("Final")
			//fmt.Println(time.Now().UnixNano() - start)

			return
		case "/api/authorize":
			ctx.SetContentType("application/json")
			code, data := handler.Authorize(ctx)
			ctx.SetStatusCode(code)
			err := json.NewEncoder(ctx).Encode(data)
			if err != nil {
				return
			}
			return
		case "/api/token":
			ctx.SetContentType("application/json")
			code, data := handler.Token(ctx)
			ctx.SetStatusCode(code)
			err := json.NewEncoder(ctx).Encode(data)
			if err != nil {
				return
			}
			return
		case "/api/me":
			ctx.SetContentType("application/json")
			code, data := handler.GetMe(ctx)
			ctx.SetStatusCode(code)
			err := json.NewEncoder(ctx).Encode(data)
			if err != nil {
				return
			}
			return
		case "/api/logout":
			ctx.SetContentType("application/json")
			code, data := handler.Logout(ctx)
			ctx.SetStatusCode(code)
			err := json.NewEncoder(ctx).Encode(data)
			if err != nil {
				return
			}
			return
		case "/api/steam":
			ctx.SetContentType("application/json")
			code, data := handler.Steam(ctx)
			ctx.SetStatusCode(code)
			err := json.NewEncoder(ctx).Encode(data)
			if err != nil {
				return
			}
			return
		case "/api/register":
			ctx.SetContentType("application/json")
			code, data := handler.Register(ctx)
			ctx.SetStatusCode(code)
			err := json.NewEncoder(ctx).Encode(data)
			if err != nil {
				return
			}
			return
		}
		path := string(ctx.Path())

		if isDevMode {
			req := &ctx.Request
			resp := &ctx.Response

			req.Header.Del("Host")
			req.SetRequestURI("http://localhost:3000" + path)

			if err := proxyClient.Do(req, resp); err != nil {
				ctx.Error("Ошибка проксирования: "+err.Error(), fasthttp.StatusBadGateway)
				return
			}
			return
		}

		filePath := "../frontend/build" + path
		if _, err := os.Stat(filePath); err == nil {
			fmt.Println("Serving file: " + path)
			fsHandler(ctx)
		} else if len(path) >= 7 && path[:7] == "/static" {
			fmt.Println("Serving static: " + path)
			fsHandler(ctx)
		} else {
			fmt.Println("Serving index.html for: " + path)
			ctx.URI().SetPath("/")
			fsHandler(ctx)
		}

		fmt.Println("static " + string(ctx.Path()))
	}

	err := fasthttp.ListenAndServe(":8080", requestHandler)
	return err
}
