package controller

import (
	"encoding/json"
	"net/http"
)

type successResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (s *successResponse) Send(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "Application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(s)

	if err != nil {
		panic(err)
	}

}

type StaffController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}
