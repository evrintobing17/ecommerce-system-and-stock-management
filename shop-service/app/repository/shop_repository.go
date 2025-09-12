package repository

import (
	"errors"
	shop "github.com/evrintobing17/ecommerce-system/shop-service/app"
	"github.com/evrintobing17/ecommerce-system/shop-service/app/models"

	"gorm.io/gorm"
)

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) shop.ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) Create(shop *models.Shop) error {
	return r.db.Create(shop).Error
}

func (r *shopRepository) FindByID(id int) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.First(&shop, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shop not found")
		}
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepository) FindByOwnerID(ownerID int, page, limit int) ([]*models.Shop, int64, error) {
	var shops []*models.Shop
	var total int64

	// Get total count
	err := r.db.Model(&models.Shop{}).Where("owner_id = ?", ownerID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err = r.db.Where("owner_id = ?", ownerID).Offset(offset).Limit(limit).Find(&shops).Error
	if err != nil {
		return nil, 0, err
	}

	return shops, total, nil
}

func (r *shopRepository) Update(shop *models.Shop) error {
	return r.db.Save(shop).Error
}

func (r *shopRepository) Delete(id int) error {
	return r.db.Delete(&models.Shop{}, "id = ?", id).Error
}
