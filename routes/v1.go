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
	productSerice := service.NewProductService(pool, productRepository)
	productController := controller.NewProductController(productSerice, validate)

	r.Handle("POST /product", Auth(http.HandlerFunc(productController.Create)))
	r.Handle("GET /product", Auth(http.HandlerFunc(productController.GetAll)))

	return r
}
