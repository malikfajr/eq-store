package entity

import "time"

type Product struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	SKU         string     `json:"sku"`
	Category    string     `json:"category"`
	ImageUrl    string     `json:"imageUrl" db:"image_url"`
	Notes       string     `json:"notes"`
	Price       int        `json:"price"`
	Stock       int        `json:"stock"`
	Location    string     `json:"location"`
	IsAvailable bool       `json:"isAvailable" db:"is_available"`
	CreatedAt   *time.Time `json:"createdAt" db:"created_at"`
}

type ProductSKU struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	SKU       string     `json:"sku"`
	Category  string     `json:"category"`
	ImageUrl  string     `json:"imageUrl" db:"image_url"`
	Price     int        `json:"price"`
	Stock     int        `json:"stock"`
	Location  string     `json:"location"`
	CreatedAt *time.Time `json:"createdAt" db:"created_at"`
}

type ProductInsertRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=30"`
	SKU         string `json:"sku" validate:"required,min=1,max=30"`
	Category    string `json:"category" validate:"required,oneof=Clothing Accessories Footwear Beverages"`
	ImageUrl    string `json:"imageUrl" validate:"required,IsURL"`
	Notes       string `json:"notes" validate:"required,min=1,max=200"`
	Price       int    `json:"price" validate:"required,min=1"`
	Stock       *int   `json:"stock" validate:"required,min=0,max=100000"`
	Location    string `json:"location" validate:"required,min=1,max=200"`
	IsAvailable *bool  `json:"isAvailable" validate:"required"`
}

type ProductUpdateRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=30"`
	SKU         string `json:"sku" validate:"required,min=1,max=30"`
	Category    string `json:"category" validate:"required,oneof=Clothing Accessories Footwear Beverages"`
	ImageUrl    string `json:"imageUrl" validate:"required,IsURL"`
	Notes       string `json:"notes" validate:"required,min=1,max=200"`
	Price       int    `json:"price" validate:"required,min=1"`
	Stock       *int   `json:"stock" validate:"required,min=1,max=100000"`
	Location    string `json:"location" validate:"required,min=1,max=200"`
	IsAvailable *bool  `json:"isAvailable" validate:"required"`
}

type ProductQueryParams struct {
	IsAvailable *bool
	InStock     *bool
	Limit       int
	Offset      int
	ID          string
	Name        string
	Category    string
	SKU         string
	Price       string
	CreatedAt   string
}
