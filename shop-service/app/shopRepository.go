package app

import (
	"github.com/evrintobing17/ecommerce-system/shop-service/app/models"
)

type ShopRepository interface {
	Create(shop *models.Shop) error
	FindByID(id int) (*models.Shop, error)
	FindByOwnerID(ownerID int, page, limit int) ([]*models.Shop, int64, error)
	Update(shop *models.Shop) error
	Delete(id int) error
}
