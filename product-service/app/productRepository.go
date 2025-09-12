package app

import "github.com/evrintobing17/ecommerce-system/product-service/app/models"

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id int) (*models.Product, error)
	FindAll(shopID int, page, limit int) ([]*models.Product, int64, error)
	Update(product *models.Product) error
	UpdateStock(id int, stock int32) error
	Delete(id int) error
}
