package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/service"
	"github.com/nyaruka/phonenumbers"
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

	if s.service.PhoneIsExist(r.Context(), body.PhoneNumber) == true {
		e := exception.NewConflict("phone is registered")
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
			return
		}
		panic(err)
	}

	success := &successResponse{
		Message: "User login successfully",
		Data:    data,
	}

	success.Send(w, http.StatusOK)
	return
}

func (s *staffController) isValidPhoneNumber(phone string) bool {
	num, err := phonenumbers.Parse(phone, "")
	if err != nil {
		log.Println("Error parsing phone number:", err)
		return false
	}

	// Memeriksa apakah nomor telepon valid
	return phonenumbers.IsPossibleNumber(num)
}
