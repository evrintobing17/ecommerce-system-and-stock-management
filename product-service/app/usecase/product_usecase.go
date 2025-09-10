package usecase

import (
	"context"
	"shared"

	product "product-service/app"
	"product-service/app/models"
)

type productUsecase struct {
	productRepo product.ProductRepository
	jwtSecret   string
	log         shared.Log
}

func NewproductUsecase(productRepo product.ProductRepository, jwtSecret string, log shared.Log) product.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
		jwtSecret:   jwtSecret,
		log:         log,
	}
}

func (s *productUsecase) ProductList(ctx context.Context) ([]models.Product, error) {
	product, err := s.productRepo.GetProductList(ctx)
	if err != nil {
		s.log.ErrorLog(err)
		return nil, err
	}

	return product, nil
}
