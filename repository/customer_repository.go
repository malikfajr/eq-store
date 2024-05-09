package repository

import (
	"context"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
)

type CustomerRepository interface {
	FindMany(ctx context.Context, pool *pgxpool.Pool, params *entity.CustomerQueryParams) *[]entity.Customer
	Create(ctx context.Context, pool *pgxpool.Pool, customer *entity.Customer) (string, error)
	IsExist(ctx context.Context, pool *pgxpool.Pool, customerId string) bool
}

type customerRepository struct{}

func NewCustomerRepository() CustomerRepository {
	return &customerRepository{}
}

func (c *customerRepository) Create(ctx context.Context, pool *pgxpool.Pool, customer *entity.Customer) (string, error) {
	query := "INSERT INTO customers (phone_number, name) VALUES ($1, $2) RETURNING id"

	err := pool.QueryRow(ctx, query, customer.PhoneNumber, customer.Name).Scan(&customer.UserId)

	return customer.UserId, err
}

func (c *customerRepository) FindMany(ctx context.Context, pool *pgxpool.Pool, params *entity.CustomerQueryParams) *[]entity.Customer {
	query := "SELECT id, phone_number, name FROM customers WHERE 1=1"
	args := pgx.NamedArgs{}

	if params.Name != "" {
		query += " AND LOWER(name) = @name"
		args["name"] = "%" + strings.ToLower(params.Name) + "%"
	}

	if params.PhoneNumber != "" {
		query += " AND phone_number = @phoneNumber"
		args["phoneNumber"] = "%" + params.PhoneNumber
	}

	rows, err := pool.Query(ctx, query, args)
	if err != nil {
		panic(err)
	}

	customers, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Customer])
	if err != nil {
		panic(err)
	}

	return &customers
}

func (c *customerRepository) IsExist(ctx context.Context, pool *pgxpool.Pool, customerId string) bool {
	var n int
	query := "SELECT 1 FROM customers WHERE id = $1"

	err := pool.QueryRow(ctx, query, customerId).Scan(&n)
	if err != nil {
		return false
	}

	log.Println(n, "exists")

	return true
}
