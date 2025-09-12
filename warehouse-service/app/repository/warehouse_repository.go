package repository

import (
	"errors"

	"github.com/evrintobing17/ecommerce-system/warehouse-service/app"
	"github.com/evrintobing17/ecommerce-system/warehouse-service/app/models"
	"gorm.io/gorm"
)

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) app.WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) Create(warehouse *models.Warehouse) error {
	return r.db.Create(warehouse).Error
}

func (r *warehouseRepository) FindByID(id int) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	err := r.db.First(&warehouse, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("warehouse not found")
		}
		return nil, err
	}
	return &warehouse, nil
}

func (r *warehouseRepository) FindByShopID(shopID int, activeOnly bool) ([]*models.Warehouse, error) {
	var warehouses []*models.Warehouse

	query := r.db.Where("shop_id = ?", shopID)
	if activeOnly {
		query = query.Where("active = ?", true)
	}

	err := query.Find(&warehouses).Error
	if err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (r *warehouseRepository) Update(warehouse *models.Warehouse) error {
	return r.db.Save(warehouse).Error
}

func (r *warehouseRepository) Delete(id int) error {
	return r.db.Delete(&models.Warehouse{}, "id = ?", id).Error
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) app.StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) Create(stock *models.Stock) error {
	return r.db.Create(stock).Error
}

func (r *stockRepository) Find(productID, warehouseID int) (*models.Stock, error) {
	var stock models.Stock
	err := r.db.First(&stock, "product_id = ? AND warehouse_id = ?", productID, warehouseID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("stock not found")
		}
		return nil, err
	}
	return &stock, nil
}

func (r *stockRepository) FindByProduct(productID int) ([]*models.Stock, error) {
	var stocks []*models.Stock
	err := r.db.Find(&stocks, "product_id = ?", productID).Error
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *stockRepository) Update(stock *models.Stock) error {
	return r.db.Save(stock).Error
}

func (r *stockRepository) Transfer(productID, fromWarehouseID, toWarehouseID int, quantity int32) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Subtract from source warehouse
		var fromStock models.Stock
		err := tx.First(&fromStock, "product_id = ? AND warehouse_id = ?", productID, fromWarehouseID).Error
		if err != nil {
			return err
		}

		if fromStock.Quantity-fromStock.Reserved < quantity {
			return errors.New("insufficient stock in source warehouse")
		}

		fromStock.Quantity -= quantity
		err = tx.Save(&fromStock).Error
		if err != nil {
			return err
		}

		// Add to destination warehouse
		var toStock models.Stock
		err = tx.First(&toStock, "product_id = ? AND warehouse_id = ?", productID, toWarehouseID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new stock record
				toStock = models.Stock{
					ProductID:   productID,
					WarehouseID: toWarehouseID,
					Quantity:    quantity,
					Reserved:    0,
					CreatedAt:   fromStock.CreatedAt,
					UpdatedAt:   fromStock.UpdatedAt,
				}
				return tx.Create(&toStock).Error
			}
			return err
		}

		toStock.Quantity += quantity
		return tx.Save(&toStock).Error
	})
}

func (r *stockRepository) Delete(productID, warehouseID int) error {
	return r.db.Delete(&models.Stock{}, "product_id = ? AND warehouse_id = ?", productID, warehouseID).Error
}
