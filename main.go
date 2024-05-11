package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/malikfajr/eq-store/middleware"
	"github.com/malikfajr/eq-store/pkg"
	"github.com/malikfajr/eq-store/routes"
)

func main() {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_PARAMS"))

	pool := pkg.CreateConnPool(context.Background(), connStr)

	defer pool.Close()

	validate := validator.New()
	validate.RegisterValidation("valid_phone", pkg.IsValidPhoneNumber)
	validate.RegisterValidation("IsURL", pkg.ValidateURL)

	r := http.NewServeMux()

	RoutesV1 := routes.NewRoutesV1(pool, validate)
	r.Handle("/v1/", http.StripPrefix("/v1", RoutesV1))

	s := http.Server{
		Addr:    ":8080",
		Handler: middleware.Logging(r),
	}

	log.Fatal(s.ListenAndServe())
}
