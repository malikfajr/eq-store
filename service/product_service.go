package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malikfajr/eq-store/entity"
	"github.com/malikfajr/eq-store/exception"
	"github.com/malikfajr/eq-store/repository"
)

type ProductService interface {
	Create(ctx context.Context, req *entity.ProductInsertRequest) (*entity.Product, error)
	GetAll(ctx context.Context, req *entity.ProductQueryParams) (*[]entity.Product, error)
	FindSku(ctx context.Context, req *entity.ProductQueryParams) (*[]entity.ProductSKU, error)
	Update(ctx context.Context, ID string, req *entity.ProductUpdateRequest) (*entity.Product, error)
	Delete(ctx context.Context, ID string) error
}

type productService struct {
	pool              *pgxpool.Pool
	productRepository repository.ProductRepository
}

func NewProductService(pool *pgxpool.Pool, productRepo repository.ProductRepository) ProductService {
	return &productService{
		pool:              pool,
		productRepository: productRepo,
	}
}

func (p *productService) Create(ctx context.Context, req *entity.ProductInsertRequest) (*entity.Product, error) {
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

	data, err := p.productRepository.Insert(ctx, p.pool, product)
	if err != nil {
		panic(exception.NewInternalServer(err.Error()))
	}

	return data, nil
}

func (p *productService) GetAll(ctx context.Context, req *entity.ProductQueryParams) (*[]entity.Product, error) {
	products, err := p.productRepository.FindMany(ctx, p.pool, req)

	return products, err
}

func (p *productService) FindSku(ctx context.Context, req *entity.ProductQueryParams) (*[]entity.ProductSKU, error) {
	productSKU, err := p.productRepository.FindSku(ctx, p.pool, req)

	return productSKU, err
}

func (p *productService) Update(ctx context.Context, ID string, req *entity.ProductUpdateRequest) (*entity.Product, error) {
	product, err := p.productRepository.FindOne(ctx, p.pool, ID)
	if err != nil {
		e := exception.NewNotFound("ID not found")
		return nil, e
	}

	product.Name = req.Name
	product.SKU = req.SKU
	product.Category = req.Category
	product.Notes = req.Notes
	product.ImageUrl = req.ImageUrl
	product.Price = req.Price
	product.Stock = *req.Stock
	product.Location = req.Location
	product.IsAvailable = *req.IsAvailable

	tx, err := p.pool.Begin(ctx)
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

	err = p.productRepository.UpdateTx(ctx, tx, product)

	if err != nil {
		panic(exception.NewInternalServer(err.Error()))
	}

	return product, nil
}

func (p *productService) Delete(ctx context.Context, ID string) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		panic(err)
	}

	err = p.productRepository.DeleteTx(ctx, tx, ID)

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

	return err
}
