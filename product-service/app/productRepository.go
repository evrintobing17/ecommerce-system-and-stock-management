package app

import (
	"context"
	"product-service/app/models"
)

type ProductRepository interface {
	GetProductList(ctx context.Context) ([]models.Product, error)
}
