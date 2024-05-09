package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/repository"
)

type CustomerService interface {
	Create(ctx context.Context, customer *entity.CustomerInsertUpdateRequest) (*entity.Customer, error)
	FindMany(ctx context.Context, params *entity.CustomerQueryParams) *[]entity.Customer
}

type customerService struct {
	pool               *pgxpool.Pool
	customerRepository repository.CustomerRepository
}

func NewCustomerService(pool *pgxpool.Pool, service repository.CustomerRepository) CustomerService {
	return &customerService{
		pool:               pool,
		customerRepository: service,
	}
}

func (c *customerService) Create(ctx context.Context, body *entity.CustomerInsertUpdateRequest) (*entity.Customer, error) {
	customer := &entity.Customer{
		Name:        body.Name,
		PhoneNumber: body.PhoneNumber,
	}

	id, err := c.customerRepository.Create(ctx, c.pool, customer)
	if err != nil {
		return nil, exception.NewConflict("phone number already exist")
	}

	customer.UserId = id
	return customer, nil
}

func (c *customerService) FindMany(ctx context.Context, params *entity.CustomerQueryParams) *[]entity.Customer {
	customers := c.customerRepository.FindMany(ctx, c.pool, params)

	return customers
}
