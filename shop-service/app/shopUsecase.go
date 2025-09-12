package app

import "github.com/evrintobing17/ecommerce-system/shop-service/app/models"

type ShopUsecase interface {
	CreateShop(name, description string, ownerID int) (*models.Shop, error)
	GetShop(id int) (*models.Shop, error)
	GetShops(ownerID, page, limit int) ([]*models.Shop, int64, error)
	UpdateShop(shop *models.Shop) error
	DeleteShop(id int) error
}
