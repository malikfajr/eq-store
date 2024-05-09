package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/service"
)

type ProductController interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	FindSku(w http.ResponseWriter, r *http.Request)
}

type productController struct {
	service  service.ProductService
	validate *validator.Validate
}

func NewProductController(service service.ProductService, validate *validator.Validate) ProductController {
	return &productController{
		service:  service,
		validate: validate,
	}
}

func (p *productController) Create(w http.ResponseWriter, r *http.Request) {
	body := entity.ProductInsertUpdateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	if err := p.validate.Struct(body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	product, err := p.service.Create(r.Context(), &body)
	if err != nil {
		e, ok := err.(*exception.CustomError)
		if ok {
			e.Send(w)
			return
		}
		w.WriteHeader(500)
		w.Write([]byte(""))
		return
	}

	data := map[string]string{
		"id":        product.Id,
		"createdAt": product.CreatedAt.Format(time.RFC3339),
	}

	success := &successResponse{
		Message: "success",
		Data:    data,
	}

	success.Send(w, http.StatusCreated)
	return
}

func (p *productController) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParams := &entity.ProductQueryParams{}

	if id := r.URL.Query().Get("id"); id != "" {
		queryParams.ID = id
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		n, err := strconv.Atoi(limit)
		if err != nil {
			queryParams.Limit = 5
		} else {
			queryParams.Limit = n
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		n, err := strconv.Atoi(offset)
		if err != nil {
			queryParams.Offset = 0
		} else {
			queryParams.Offset = n
		}
	}

	if name := r.URL.Query().Get("name"); name != "" {
		queryParams.Name = name
	}

	if isAvailable := r.URL.Query().Get("isAvailable"); isAvailable != "" {
		available, err := strconv.ParseBool(isAvailable)
		if err != nil {
			queryParams.IsAvailable = nil
		} else {
			queryParams.IsAvailable = &available
		}
	}

	if category := r.URL.Query().Get("category"); category != "" {
		if p.isValidCategory(category) {
			queryParams.Category = category
		}
	}

	if sku := r.URL.Query().Get("sku"); sku != "" {
		queryParams.SKU = sku
	}

	if inStock := r.URL.Query().Get("inStock"); inStock != "" {
		stock, err := strconv.ParseBool(inStock)
		if err != nil {
			queryParams.InStock = nil
		} else {
			queryParams.InStock = &stock
		}
	}

	if price := r.URL.Query().Get("price"); price != "" {
		if p.isValidOrder(price) {
			queryParams.Price = price
		}
	}

	if createdAt := r.URL.Query().Get("createdAt"); createdAt != "" {
		if p.isValidOrder(createdAt) {
			queryParams.CreatedAt = createdAt
		}
	}

	data, err := p.service.GetAll(r.Context(), queryParams)
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
		Data:    data,
	}

	success.Send(w, http.StatusOK)
	return
}

func (p *productController) Update(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("id")
	body := entity.ProductInsertUpdateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	if err := p.validate.Struct(body); err != nil {
		e := exception.NewBadRequest("request doesn’t pass validation")
		e.Send(w)
		return
	}

	product, err := p.service.Update(r.Context(), ID, &body)
	if err != nil {
		e, ok := err.(*exception.CustomError)
		if ok {
			e.Send(w)
			return
		}
		panic(e)
	}

	success := &successResponse{
		Message: "Success update product",
		Data:    product,
	}

	success.Send(w, http.StatusOK)
	return
}

func (p *productController) Delete(w http.ResponseWriter, r *http.Request) {
	ID := r.PathValue("id")

	err := p.service.Delete(r.Context(), ID)
	if err != nil {
		e, ok := err.(*exception.CustomError)
		if ok {
			e.Send(w)
			return
		}
		panic(e)
	}

	success := &successResponse{
		Message: "Delete product success",
		Data:    []string{},
	}

	success.Send(w, http.StatusOK)
	return
}

func (p *productController) FindSku(w http.ResponseWriter, r *http.Request) {
	queryParams := &entity.ProductQueryParams{}

	if id := r.URL.Query().Get("id"); id != "" {
		queryParams.ID = id
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		n, err := strconv.Atoi(limit)
		if err != nil {
			queryParams.Limit = 5
		} else {
			queryParams.Limit = n
		}
	}

	if offset := r.URL.Query().Get("offset"); offset != "" {
		n, err := strconv.Atoi(offset)
		if err != nil {
			queryParams.Offset = 0
		} else {
			queryParams.Offset = n
		}
	}

	if name := r.URL.Query().Get("name"); name != "" {
		queryParams.Name = name
	}

	if category := r.URL.Query().Get("category"); category != "" {
		if p.isValidCategory(category) {
			queryParams.Category = category
		}
	}

	if sku := r.URL.Query().Get("sku"); sku != "" {
		queryParams.SKU = sku
	}

	if inStock := r.URL.Query().Get("inStock"); inStock != "" {
		stock, err := strconv.ParseBool(inStock)
		if err != nil {
			queryParams.InStock = nil
		} else {
			queryParams.InStock = &stock
		}
	}

	if price := r.URL.Query().Get("price"); price != "" {
		if p.isValidOrder(price) {
			queryParams.Price = price
		}
	}

	data, err := p.service.FindSku(r.Context(), queryParams)
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
		Data:    data,
	}

	success.Send(w, http.StatusOK)
	return
}

func (p *productController) isValidCategory(key string) bool {
	ok := false
	categories := map[string]bool{
		"Clothing":    true,
		"Accessories": true,
		"Footwear":    true,
		"Beverages":   true,
	}

	_, ok = categories[key]
	return ok
}

func (p *productController) isValidOrder(key string) bool {
	ok := false
	order := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	_, ok = order[key]
	return ok
}
