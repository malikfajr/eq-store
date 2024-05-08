package exception

import (
	"encoding/json"
	"net/http"
)

type CustomError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func (c *CustomError) Send(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "Application/json")
	w.WriteHeader(c.StatusCode)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(c)

	if err != nil {
		panic(err)
	}
}

func (c *CustomError) Error() string {
	return c.Message
}

func (c *CustomError) Status() int {
	return c.StatusCode
}

func NewBadRequest(message string) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

func NewNotFound(message string) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

func NewConflict(message string) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

func NewUnauthorized(message string) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

func NewInternalServer(message string) *CustomError {
	return &CustomError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}
