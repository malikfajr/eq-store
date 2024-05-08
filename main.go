package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/middleware"
	"github.com/malikfajr/eq-store/routes"
)

func main() {

	pool, err := pgxpool.New(context.Background(), "postgres://postgres:secret@localhost:5432/eq?sslmode=disable")
	if err != nil {
		log.Fatalln("Cannot connect database: ", err)
		os.Exit(1)
	}
	log.Println("Database connected")
	defer pool.Close()

	validate := validator.New()

	r := http.NewServeMux()

	RoutesV1 := routes.NewRoutesV1(pool, validate)
	r.Handle("/v1/", http.StripPrefix("/v1", RoutesV1))

	s := http.Server{
		Addr:    ":8080",
		Handler: middleware.Logging(r),
	}

	log.Fatal(s.ListenAndServe())
}
