package service

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/repository"
)

type ProductService interface {
	Create(ctx context.Context, req entity.ProductInsertUpdateRequest) (*entity.Product, error)
	GetAll(ctx context.Context, req entity.ProductQueryParams) (*[]entity.Product, error)
}

type productService struct {
	pool               *pgxpool.Pool
	productRepositoory repository.ProductRepository
}

func NewProductService(pool *pgxpool.Pool, productRepo repository.ProductRepository) ProductService {
	return &productService{
		pool:               pool,
		productRepositoory: productRepo,
	}
}

func (p *productService) Create(ctx context.Context, req entity.ProductInsertUpdateRequest) (*entity.Product, error) {
	product := &entity.Product{
		Name:        req.Name,
		SKU:         req.SKU,
		Category:    req.Category,
		ImageUrl:    req.ImageUrl,
		Notes:       req.Notes,
		Price:       req.Price,
		Stock:       *req.Stock,
		Location:    req.Location,
		IsAvailable: *req.IsAvailable,
	}

	data, err := p.productRepositoory.Insert(ctx, p.pool, product)
	if err != nil {
		log.Println(err)
		e := exception.NewInternalServer("Internal server error")
		return nil, e
	}

	return data, nil
}

func (p *productService) GetAll(ctx context.Context, req entity.ProductQueryParams) (*[]entity.Product, error) {
	products, err := p.productRepositoory.FindMany(ctx, p.pool, &req)

	return products, err
}
