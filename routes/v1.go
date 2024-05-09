package routes

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/controller"
	"github.com/malikfajr/eq-store/middleware"
	"github.com/malikfajr/eq-store/repository"
	"github.com/malikfajr/eq-store/service"
)

func NewRoutesV1(pool *pgxpool.Pool, validate *validator.Validate) *http.ServeMux {
	Auth := middleware.Auth

	r := http.NewServeMux()

	staffRepository := repository.NewStaffRepository()
	staffService := service.NewStaffService(staffRepository, pool)
	staffController := controller.NewStaffController(staffService, validate)

	r.HandleFunc("POST /staff/register", staffController.Register)
	r.HandleFunc("POST /staff/login", staffController.Login)

	productRepository := repository.NewProductRepository()
	productService := service.NewProductService(pool, productRepository)
	productController := controller.NewProductController(productService, validate)

	r.Handle("POST /product", Auth(http.HandlerFunc(productController.Create)))
	r.Handle("GET /product", Auth(http.HandlerFunc(productController.GetAll)))
	r.Handle("PUT /product/{id}", Auth(http.HandlerFunc(productController.Update)))
	r.Handle("DELETE /product/{id}", Auth(http.HandlerFunc(productController.Delete)))

	r.Handle("GET /product/customer", http.HandlerFunc(productController.FindSku))

	customerRepoitory := repository.NewCustomerRepository()
	customerService := service.NewCustomerService(pool, customerRepoitory)
	customerController := controller.NewCustomerController(validate, customerService)

	r.Handle("POST /customer/register", Auth(http.HandlerFunc(customerController.Create)))
	r.Handle("GET /customer", Auth(http.HandlerFunc(customerController.GetAll)))

	transactionRepository := repository.NewTransactionRepository()
	transactionService := service.NewTransactionService(pool, customerRepoitory, productRepository, transactionRepository)
	transactionController := controller.NewTransactionController(validate, transactionService)

	r.Handle("POST /product/checkout", Auth(http.HandlerFunc(transactionController.Create)))
	r.Handle("GET /product/checkout/history", Auth(http.HandlerFunc(transactionController.GetAll)))
	return r
}
