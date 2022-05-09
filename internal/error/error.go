package apierror

import (
	"errors"
)

var LowBalance = errors.New("low balance")

type ExtCode int

const (
	NotEnoughMoney ExtCode = 1
	BindingError   ExtCode = 2
	NotFound       ExtCode = 3
)

type APIError struct {
	Message string `json:"message"`
	// Расширенный код ошибки
	ExtCode ExtCode `json:"ext_code,omitempty"`
}
