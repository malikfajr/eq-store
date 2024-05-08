package routes

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/controller"
	"github.com/malikfajr/eq-store/repository"
	"github.com/malikfajr/eq-store/service"
)

func NewStaffRoute(pool *pgxpool.Pool, validate *validator.Validate) *http.ServeMux {
	staffRepository := repository.NewStaffRepository()
	staffService := service.NewStaffService(staffRepository, pool)
	staffController := controller.NewStaffController(staffService, validate)

	r := http.NewServeMux()

	r.HandleFunc("POST /register", staffController.Register)
	r.HandleFunc("POST /login", staffController.Login)

	return r
}
