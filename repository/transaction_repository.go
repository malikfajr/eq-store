package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx pgx.Tx, payload *entity.TransactionInsertRequest) string
	InsertDetail(ctx context.Context, tx pgx.Tx, transactionId string, payload []entity.ProductDetail)
	DecrementStock(ctx context.Context, tx pgx.Tx, payload []entity.ProductDetail)
	FindMany(ctx context.Context, pool *pgxpool.Pool, payload *entity.TransactionQueryParams) []entity.Transaction
}

type transactionRepository struct{}

func NewTransactionRepository() TransactionRepository {
	return &transactionRepository{}
}

func (t *transactionRepository) Create(ctx context.Context, tx pgx.Tx, payload *entity.TransactionInsertRequest) string {
	var id string
	query := "INSERT INTO transactions (customer_id, paid, change) VALUES ($1, $2, $3) RETURNING id"

	err := tx.QueryRow(ctx, query, payload.CustomerId, payload.Paid, payload.Change).Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}

func (t *transactionRepository) InsertDetail(ctx context.Context, tx pgx.Tx, transactionId string, payload []entity.ProductDetail) {
	tx.CopyFrom(
		ctx,
		pgx.Identifier{"transaction_detail"},
		[]string{"transaction_id", "product_id", "quantity"},
		pgx.CopyFromSlice(len(payload), func(i int) ([]interface{}, error) {
			return []interface{}{transactionId, payload[i].ProductId, payload[i].Quantity}, nil
		}),
	)
}

func (t *transactionRepository) DecrementStock(ctx context.Context, tx pgx.Tx, payload []entity.ProductDetail) {
	query := "UPDATE products SET stock = stock - $1 WHERE id = $2"

	for _, pd := range payload {
		_, err := tx.Exec(ctx, query, pd.Quantity, pd.ProductId)
		if err != nil {
			panic(err)
		}
	}
}

func (t *transactionRepository) FindMany(ctx context.Context, pool *pgxpool.Pool, params *entity.TransactionQueryParams) []entity.Transaction {
	query := `
		SELECT t.id, t.customer_id, t.paid, t.change, t.created_at, 
			(SELECT JSON_AGG(json_build_object('productId', td.product_id, 'quantity', td.quantity)) 
				FROM transaction_detail td 
				WHERE td.transaction_id = t.id) AS pd_details 
		FROM transactions AS t WHERE 1=1`

	args := pgx.NamedArgs{}

	if params.CustomerId != "" {
		query += " AND t.customer_id = @customerId"
		args["customerId"] = params.CustomerId
	}

	if params.CreatedAt != "" {
		query += " ORDER BY t.created_at " + params.CreatedAt
	} else {
		query += " ORDER BY t.created_at desc"
	}

	query += " LIMIT @limit OFFSET @offset"
	args["limit"] = params.Limit
	args["offset"] = params.Offset

	rows, err := pool.Query(ctx, query, args)
	if err != nil {
		panic(err)
	}

	var transactions []entity.Transaction = []entity.Transaction{}
	for rows.Next() {
		transaction := &entity.Transaction{}
		rows.Scan(&transaction.Id, &transaction.CustomerId, &transaction.Paid, &transaction.Change, &transaction.CreatedAt, &transaction.ProductDetails)
		transactions = append(transactions, *transaction)
	}

	return transactions
}
