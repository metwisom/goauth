package responseError

const (
	ErrMethodNotAllowed    = "Метод не поддерживается"
	ErrMissingParameters   = "Missing required parameters"
	ErrMissingCookie       = "Session cookie is missing"
	ErrInvalidSession      = "Invalid session"
	ErrClientIdError       = "client_id must be an integer"
	ErrInvalidRedirectUri  = "Invalid redirect_uri format"
	ErrInvalidCredentials  = "Неверный логин или пароль"
	ErrSessionCreationFail = "Ошибка при создании сессии"
	ErrMissingCredentials  = "Логин и пароль обязательны"
)

type ErrorResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

func BadRequest(message string) ErrorResponse {
	return ErrorResponse{
		Ok:          false,
		ErrorCode:   400,
		Description: "Bad Request: " + message,
	}
}

func Unauthorized(message string) ErrorResponse {
	return ErrorResponse{
		Ok:          false,
		ErrorCode:   401,
		Description: "Unauthorized: " + message,
	}
}

func InternalServerError(message string) ErrorResponse {
	return ErrorResponse{
		Ok:          false,
		ErrorCode:   500,
		Description: "Internal Server Error: " + message,
	}
}

func MethodNotAllowed(message string) ErrorResponse {
	return ErrorResponse{
		Ok:          false,
		ErrorCode:   405,
		Description: "Method Not Allowed: " + message,
	}
}
