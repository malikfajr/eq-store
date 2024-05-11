package entity

import "time"

type ProductDetail struct {
	TransactionId string `json:"-"`
	ProductId     string `json:"productId" validate:"required"`
	Quantity      int    `json:"quantity" validate:"required,min=1"`
	Price         int    `json:"-"`
	TotalPrice    int    `json:"-"`
}

type Transaction struct {
	Id             string          `json:"transactionId"`
	CustomerId     string          `json:"customerId"`
	Paid           int             `json:"paid"`
	Change         int             `json:"change"`
	ProductDetails []ProductDetail `json:"productDetails"`
	CreatedAt      *time.Time      `json:"createdAt" db:"created_at"`
}

type TransactionInsertRequest struct {
	CustomerId     string          `json:"customerId" validate:"required"`
	ProductDetails []ProductDetail `json:"productDetails" validate:"required,gte=1,dive,required"` // TODO: validate if product id duplicate fi
	Paid           int             `json:"paid" validate:"required,min=1"`
	Change         *int            `json:"change" validate:"required,min=0"`
}

type TransactionQueryParams struct {
	Limit      int
	Offset     int
	CustomerId string
	CreatedAt  string
}
