package repository

import (
	"errors"

	product "github.com/evrintobing17/ecommerce-system/product-service/app"
	"github.com/evrintobing17/ecommerce-system/product-service/app/models"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) product.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id int) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindAll(shopID, page, limit int) ([]*models.Product, int64, error) {
	var products []*models.Product
	var total int64

	query := r.db.Model(&models.Product{})
	if shopID != 0 {
		query = query.Where("shop_id = ?", shopID)
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) UpdateStock(id int, stock int32) error {
	return r.db.Model(&models.Product{}).Where("id = ?", id).Update("stock", stock).Error
}

func (r *productRepository) Delete(id int) error {
	return r.db.Delete(&models.Product{}, "id = ?", id).Error
}
