package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/service"
)

type TransactionController interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
}

type transactionController struct {
	transactionService service.TransactionService
	validate           *validator.Validate
}

func NewTransactionController(validate *validator.Validate, transactionSerive service.TransactionService) TransactionController {
	return &transactionController{
		validate:           validate,
		transactionService: transactionSerive,
	}
}

// Create implements TransactionController.
func (t *transactionController) Create(w http.ResponseWriter, r *http.Request) {
	body := &entity.TransactionInsertRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	if err := t.validate.Struct(body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	err := t.transactionService.Create(r.Context(), body)
	if err != nil {
		log.Println(err)
		e, ok := err.(*exception.CustomError)
		if ok {
			e.Send(w)
			return
		}
		panic(err)
	}

	success := &successResponse{
		Message: "success",
		Data:    make([]string, 0),
	}

	success.Send(w, http.StatusCreated)
}

// GetAll implements TransactionController.
func (t *transactionController) GetAll(w http.ResponseWriter, r *http.Request) {
	params := &entity.TransactionQueryParams{}

	if customerId := r.URL.Query().Get("customerId"); customerId != "" {
		params.CustomerId = customerId
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		n, err := strconv.Atoi(limit)
		if err != nil {
			params.Limit = 5
		} else {
			params.Limit = n
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		n, err := strconv.Atoi(offset)
		if err != nil {
			params.Offset = 0
		} else {
			params.Offset = n
		}
	}

	if createdAt := r.URL.Query().Get("createdAt"); createdAt != "" {
		if t.isValidOrder(createdAt) {
			params.CreatedAt = createdAt
		}
	}

	data, err := t.transactionService.FindMany(r.Context(), params)
	if err != nil {
		e, ok := err.(*exception.CustomError)
		if ok {
			e.Send(w)
			return
		}
		panic(err)
	}

	success := &successResponse{
		Message: "success",
		Data:    data,
	}

	success.Send(w, http.StatusOK)
}

func (t *transactionController) isValidateInsertPayload(payload *entity.TransactionInsertRequest) error {
	if err := t.validate.Struct(payload); err != nil {
		return exception.NewBadRequest("request doesn’t pass validation")
	}

	return nil
}

func (t *transactionController) isValidOrder(key string) bool {
	order := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	_, ok := order[key]

	return ok
}
