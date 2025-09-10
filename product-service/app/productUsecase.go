package app

import (
	"context"
	"product-service/app/models"
)

type ProductUsecase interface {
	ProductList(ctx context.Context) ([]models.Product, error)
}
