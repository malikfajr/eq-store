package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/repository"
)

type TransactionService interface {
	Create(ctx context.Context, payload *entity.TransactionInsertRequest) error
	FindMany(ctx context.Context, params *entity.TransactionQueryParams) (*[]entity.Transaction, error)
}

type transactionService struct {
	pool                  *pgxpool.Pool
	customerRepository    repository.CustomerRepository
	productRepository     repository.ProductRepository
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(pool *pgxpool.Pool, customerRepository repository.CustomerRepository, productRepository repository.ProductRepository, transactionRepository repository.TransactionRepository) TransactionService {
	return &transactionService{
		pool:                  pool,
		customerRepository:    customerRepository,
		productRepository:     productRepository,
		transactionRepository: transactionRepository,
	}
}

func (t *transactionService) Create(ctx context.Context, payload *entity.TransactionInsertRequest) error {
	if err := t.isValidPayload(ctx, payload); err != nil {
		return err
	}

	tx, err := t.pool.Begin(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := recover()
		if err != nil {
			e := tx.Rollback(ctx)
			if e != nil {
				panic(e)
			}
		} else {
			e := tx.Commit(ctx)
			if e != nil {
				panic(e)
			}
		}
	}()

	id := t.transactionRepository.Create(ctx, tx, payload)
	t.transactionRepository.InsertDetail(ctx, tx, id, payload.ProductDetails)
	t.transactionRepository.DecrementStock(ctx, tx, payload.ProductDetails)

	return nil
}

func (t *transactionService) FindMany(ctx context.Context, params *entity.TransactionQueryParams) (*[]entity.Transaction, error) {
	transactions := t.transactionRepository.FindMany(ctx, t.pool, params)
	return &transactions, nil
}

func (t *transactionService) isValidPayload(ctx context.Context, payload *entity.TransactionInsertRequest) error {
	// 0. customer id exists - 404
	if exists := t.customerRepository.IsExist(ctx, t.pool, payload.CustomerId); !exists {
		return exception.NewNotFound("Customer id not found")
	}

	// 1. product id exists - 404
	var productDetails map[string]int = map[string]int{}
	productIds := []string{}

	for _, product := range payload.ProductDetails {
		productDetails[product.ProductId] = product.Quantity
		productIds = append(productIds, product.ProductId)
	}

	products := t.productRepository.FindByIds(ctx, t.pool, productIds)
	if len(*products) != len(productDetails) {
		return exception.NewNotFound("one of productId not found")
	}

	// 2. paid is enought - 400
	totalPrice := 0

	for _, product := range *products {
		if product.IsAvailable == false { // 5. one of product isAvailable false - 400
			return exception.NewBadRequest("one of product not available")
		}

		if product.Stock < productDetails[product.Id] { // 4. product stock is enought - 400
			return exception.NewBadRequest("one of productIds stock is not enough")
		}
		totalPrice += (product.Price * productDetails[product.Id])
	}

	if totalPrice > payload.Paid {
		return exception.NewBadRequest("paid is not enough based on all bought product")
	}

	// 3. change is right - 400
	if change := payload.Paid - totalPrice; change != *payload.Change {
		return exception.NewBadRequest("change is not right")
	}

	return nil
}
