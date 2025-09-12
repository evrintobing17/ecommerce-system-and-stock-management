package app

import "github.com/evrintobing17/ecommerce-system/product-service/app/models"

type ProductUsecase interface {
	GetProducts(shopID int, page, limit int) ([]*models.Product, int64, error)
	GetProduct(id int) (*models.Product, error)
	CreateProduct(name, description string, price float64, stock int32, shopID int) (*models.Product, error)
	UpdateProduct(product *models.Product) error
	AddStock(id int, quantity int32) error
	SubtractStock(id int, quantity int32) error
	SetStock(id int, quantity int32) error
	DeleteProduct(id int) error
}
