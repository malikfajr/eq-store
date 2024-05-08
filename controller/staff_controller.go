package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/service"
)

type staffController struct {
	service  service.StaffService
	validate *validator.Validate
}

func NewStaffController(service service.StaffService, validate *validator.Validate) StaffController {
	return &staffController{
		service:  service,
		validate: validate,
	}
}

// Register implements StaffController.
func (s *staffController) Register(w http.ResponseWriter, r *http.Request) {
	body := &entity.StaffRegisterRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	if err := s.validate.Struct(body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	data, err := s.service.Register(r.Context(), body)
	if err != nil {
		e := exception.NewConflict("phone number is exists")
		e.Send(w)
		return
	}

	success := &successResponse{
		Message: "User registered successfully",
		Data:    data,
	}

	success.Send(w, http.StatusCreated)
	return

}

func (s *staffController) Login(w http.ResponseWriter, r *http.Request) {
	body := &entity.StaffLoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	if err := s.validate.Struct(body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	data, err := s.service.Login(r.Context(), body)
	if err != nil {
		if e, ok := err.(*exception.CustomError); ok {
			e.Send(w)
		}
		return
	}

	success := &successResponse{
		Message: "User login successfully",
		Data:    data,
	}

	success.Send(w, http.StatusOK)
	return
}
