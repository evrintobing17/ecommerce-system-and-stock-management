package usecase

import (
	"time"

	"github.com/evrintobing17/ecommerce-system/shop-service/app"
	"github.com/evrintobing17/ecommerce-system/shop-service/app/models"
)

type shopUsecase struct {
	shopRepo app.ShopRepository
}

func NewShopUsecase(shopRepo app.ShopRepository) app.ShopUsecase {
	return &shopUsecase{shopRepo: shopRepo}
}

func (u *shopUsecase) CreateShop(name, description string, ownerID int) (*models.Shop, error) {
	shop := &models.Shop{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := u.shopRepo.Create(shop)
	if err != nil {
		return nil, err
	}

	return &models.Shop{
		ID:          shop.ID,
		Name:        shop.Name,
		Description: shop.Description,
		OwnerID:     shop.OwnerID,
		CreatedAt:   shop.CreatedAt,
		UpdatedAt:   shop.UpdatedAt,
	}, nil
}

func (u *shopUsecase) GetShop(id int) (*models.Shop, error) {
	shop, err := u.shopRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &models.Shop{
		ID:          shop.ID,
		Name:        shop.Name,
		Description: shop.Description,
		OwnerID:     shop.OwnerID,
		CreatedAt:   shop.CreatedAt,
		UpdatedAt:   shop.UpdatedAt,
	}, nil
}

func (u *shopUsecase) GetShops(ownerID, page, limit int) ([]*models.Shop, int64, error) {
	shops, total, err := u.shopRepo.FindByOwnerID(ownerID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []*models.Shop
	for _, shop := range shops {
		result = append(result, &models.Shop{
			ID:          shop.ID,
			Name:        shop.Name,
			Description: shop.Description,
			OwnerID:     shop.OwnerID,
			CreatedAt:   shop.CreatedAt,
			UpdatedAt:   shop.UpdatedAt,
		})
	}

	return result, total, nil
}

func (u *shopUsecase) UpdateShop(shop *models.Shop) error {
	existingShop, err := u.shopRepo.FindByID(shop.ID)
	if err != nil {
		return err
	}

	existingShop.Name = shop.Name
	existingShop.Description = shop.Description
	existingShop.UpdatedAt = time.Now()

	return u.shopRepo.Update(existingShop)
}

func (u *shopUsecase) DeleteShop(id int) error {
	return u.shopRepo.Delete(id)
}
