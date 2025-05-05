package utils

import (
	"errors"
	"fmt"
)

type MyError struct {
	Code    int
	Message string
	Context string
}

// Реализация интерфейса error
func (e *MyError) Error() string {
	return fmt.Sprintf("Error %d: %s (context: %s)", e.Code, e.Message, e.Context)
}

// Is Метод для сравнения ошибок
func (e *MyError) Is(target error) bool {
	var t *MyError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	// Сравниваем только Code для упрощения, можно добавить другие поля
	return e.Code == t.Code
}

// NewMyError Конструктор для создания ошибок
func NewMyError(code int, message, context string) *MyError {
	return &MyError{
		Code:    code,
		Message: message,
		Context: context,
	}
}
