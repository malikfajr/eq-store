package repository

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
)

type ProductRepository interface {
	Insert(ctx context.Context, pool *pgxpool.Pool, product *entity.Product) (*entity.Product, error)
	IsExists(ctx context.Context, pool *pgxpool.Pool, productId string) bool
	FindMany(ctx context.Context, pool *pgxpool.Pool, params *entity.ProductQueryParams) (*[]entity.Product, error)
	FindSku(ctx context.Context, pool *pgxpool.Pool, params *entity.ProductQueryParams) (*[]entity.ProductSKU, error)
	FindOne(ctx context.Context, pool *pgxpool.Pool, ID string) (*entity.Product, error)
	FindByIds(ctx context.Context, pool *pgxpool.Pool, productIds []string) *[]entity.Product
	UpdateTx(ctx context.Context, tx pgx.Tx, product *entity.Product) error
	Delete(ctx context.Context, pool *pgxpool.Pool, ID string) error
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
	query := "SELECT id, name, sku, category, image_url, notes, price, stock, location, is_available, created_at FROM products WHERE deleted_at IS NULL"
	args := pgx.NamedArgs{}

	if params.ID != "" {
		query += " AND id = @id"
		args["id"] = params.ID
	}

	if params.Name != "" {
		query += " AND LOWER(name) = @name"
		args["name"] = "%" + strings.ToLower(params.Name) + "%"
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
	} else {
		query += " ORDER BY  created_at desc"
	}

	args["limit"] = params.Limit
	args["offset"] = params.Offset

	query += " LIMIT @limit OFFSET @offset"

	rows, err := pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	products, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Product])

	return &products, err
}

func (p *productRepository) IsExists(ctx context.Context, pool *pgxpool.Pool, productId string) bool {
	var id int
	query := "SELECT 1 FROM products WHERE deleted_at IS NULL AND id = $1 "

	err := pool.QueryRow(ctx, query, productId).Scan(&id)
	if err != nil {
		return false
	}

	log.Println(id, "is found")
	return true
}

func (p *productRepository) FindOne(ctx context.Context, pool *pgxpool.Pool, ID string) (*entity.Product, error) {
	var product entity.Product
	query := "SELECT id, name, sku, category, image_url, notes, price, stock, location, is_available, created_at FROM products WHERE deleted_at IS NULL AND id = $1 LIMIT 1;"

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

func (p *productRepository) FindSku(ctx context.Context, pool *pgxpool.Pool, params *entity.ProductQueryParams) (*[]entity.ProductSKU, error) {
	query := "SELECT id, name, sku, category, image_url, price, stock, location, created_at FROM products WHERE deleted_at IS NULL AND is_available = true"
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

	if params.InStock != nil {
		if *params.InStock == true {
			query += " AND stock > 0"
		} else {
			query += " AND stock = 0"
		}
	}

	if params.Price != "" {
		query += " ORDER BY price " + params.Price
	} else {
		query += " ORDER BY created_at desc"
	}

	args["limit"] = params.Limit
	args["offset"] = params.Offset

	query += " LIMIT @limit OFFSET @offset"

	rows, err := pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	productSKU, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.ProductSKU])

	return &productSKU, err
}

func (p *productRepository) UpdateTx(ctx context.Context, tx pgx.Tx, product *entity.Product) error {
	query := `
		UPDATE products 
			SET name = @name, sku = @sku, category = @category, image_url = @imageUrl, 
				notes = @notes, price = @price, stock = @stock, location = @location, is_available = @isAvailable
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":          product.Id,
		"name":        product.Name,
		"sku":         product.SKU,
		"category":    product.Category,
		"imageUrl":    product.ImageUrl,
		"notes":       product.Notes,
		"price":       product.Price,
		"stock":       product.Stock,
		"location":    product.Location,
		"isAvailable": product.IsAvailable,
	}

	_, err := tx.Exec(ctx, query, args)

	return err
}

func (p *productRepository) Delete(ctx context.Context, pool *pgxpool.Pool, ID string) error {
	query := "UPDATE products SET deleted_at = NOW() WHERE ID = $1"

	tag, err := pool.Exec(ctx, query, ID)

	if tag.RowsAffected() < 1 {
		return exception.NewNotFound("product id not found")
	}

	return err
}

func (p *productRepository) FindByIds(ctx context.Context, pool *pgxpool.Pool, productIds []string) *[]entity.Product {
	query := "SELECT id, name, sku, category, image_url, notes, price, stock, location, is_available, created_at FROM products WHERE deleted_at IS NULL AND id::TEXT = ANY($1);"

	rows, err := pool.Query(ctx, query, productIds)
	if err != nil {
		panic(err)
	}

	products, err := pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Product])
	if err != nil {
		panic(err)
	}
	return &products
}
