package app

import "github.com/evrintobing17/ecommerce-system/warehouse-service/app/models"

type WarehouseRepository interface {
	Create(warehouse *models.Warehouse) error
	FindByID(id int) (*models.Warehouse, error)
	FindByShopID(shopID int, activeOnly bool) ([]*models.Warehouse, error)
	Update(warehouse *models.Warehouse) error
	Delete(id int) error
}

type StockRepository interface {
	Create(stock *models.Stock) error
	Find(productID, warehouseID int) (*models.Stock, error)
	FindByProduct(productID int) ([]*models.Stock, error)
	Update(stock *models.Stock) error
	Transfer(productID, fromWarehouseID, toWarehouseID int, quantity int32) error
	Delete(productID, warehouseID int) error
}
