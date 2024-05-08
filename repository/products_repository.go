package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
)

type ProductRepository interface {
	Insert(ctx context.Context, pool *pgxpool.Pool, product *entity.Product) (*entity.Product, error)
	FindMany(ctx context.Context, pool *pgxpool.Pool, params *entity.ProductQueryParams) (*[]entity.Product, error)
	FindOne(ctx context.Context, pool *pgxpool.Pool, ID string) (*entity.Product, error)
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (p *productRepository) Insert(ctx context.Context, pool *pgxpool.Pool, product *entity.Product) (*entity.Product, error) {
	query := `
		INSERT INTO products (name, sku, category, image_url, notes, price, stock, location, is_available)
		VALUES (@name, @sku, @category, @imageUrl, @notes, @price, @stock, @location, @isAvailable)
		RETURNING id, created_at
	`

	args := pgx.NamedArgs{
		"name":        product.Name,
		"sku":         product.SKU,
		"category":    product.Category,
		"imageUrl":    product.ImageUrl,
		"notes":       product.Notes,
		"price":       product.Price,
		"stock":       product.Price,
		"location":    product.Location,
		"isAvailable": product.IsAvailable,
	}

	err := pool.QueryRow(ctx, query, args).Scan(&product.Id, &product.CreatedAt)
	return product, err
}

func (p *productRepository) FindMany(ctx context.Context, pool *pgxpool.Pool, params *entity.ProductQueryParams) (*[]entity.Product, error) {
	query := "SELECT * FROM products WHERE 1=1"
	args := pgx.NamedArgs{}

	if params.ID != "" {
		query += " AND id = @id"
		args["id"] = params.ID
	}

	if params.Name != "" {
		query += " AND name = @name"
		args["name"] = "%" + params.Name + "%"
	}

	if params.SKU != "" {
		query += " AND sku = @sku"
		args["sku"] = params.SKU
	}

	if params.Category != "" {
		query += " AND category = @category"
		args["category"] = params.Category
	}

	if params.IsAvailable != nil {
		query += " AND is_available = @isAvailable"
		args["isAvailable"] = *params.IsAvailable
	}

	if params.InStock != nil {
		if *params.InStock == true {
			query += " AND stock > 0"
		} else {
			query += " AND stock = 0"
		}
	}

	if params.Price != "" && params.CreatedAt != "" {
		query += " ORDER BY price " + params.Price + ", created_at " + params.CreatedAt
	} else if params.Price != "" {
		query += " ORDER BY price " + params.Price
	} else if params.CreatedAt != "" {
		query += " ORDER BY created_at " + params.CreatedAt
	}

	if params.Limit <= 0 {
		params.Limit = 5
	}

	args["limit"] = params.Limit
	args["offset"] = params.Offset

	query += " LIMIT @limit OFFSET @offset"

	log.Println(args, query)

	rows, err := pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	products, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Product])

	return &products, err
}

func (p *productRepository) FindOne(ctx context.Context, pool *pgxpool.Pool, ID string) (*entity.Product, error) {
	var product entity.Product
	query := "SELECT * FROM products WHERE id = $1 LIMIT 1;"

	rows, err := pool.Query(ctx, query, ID)
	if err != nil {
		return nil, errors.New("product id not found")
	}

	product, err = pgx.CollectOneRow(rows, pgx.RowToStructByPos[entity.Product])
	if err != nil {
		return nil, errors.New("product id not found")
	}

	return &product, nil
}
