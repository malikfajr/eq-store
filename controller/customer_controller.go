package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/service"
)

type CustomerController interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
}

type customerController struct {
	customerService service.CustomerService
	validate        *validator.Validate
}

func NewCustomerController(validate *validator.Validate, service service.CustomerService) CustomerController {
	return &customerController{
		validate:        validate,
		customerService: service,
	}
}

func (c *customerController) Create(w http.ResponseWriter, r *http.Request) {
	body := &entity.CustomerInsertUpdateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	if err := c.validate.Struct(body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	customer, err := c.customerService.Create(r.Context(), body)
	if err != nil {
		e, ok := err.(*exception.CustomError)
		if ok {
			e.Send(w)
			return
		}
		panic(e)
	}

	success := &successResponse{
		Message: "Create custommer success",
		Data:    customer,
	}

	success.Send(w, http.StatusCreated)
}

func (c *customerController) GetAll(w http.ResponseWriter, r *http.Request) {
	params := &entity.CustomerQueryParams{}

	if name := r.URL.Query().Get("name"); name != "" {
		params.Name = name
	}

	if phone := r.URL.Query().Get("phoneNumber"); phone != "" {
		params.PhoneNumber = phone
	}

	customers := c.customerService.FindMany(r.Context(), params)

	success := &successResponse{
		Message: "success",
		Data:    customers,
	}

	success.Send(w, http.StatusOK)
}
