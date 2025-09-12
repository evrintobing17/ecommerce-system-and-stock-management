package app

import "github.com/evrintobing17/ecommerce-system/warehouse-service/app/models"

type WarehouseUsecase interface {
	GetWarehouse(id int) (*models.Warehouse, error)
	GetWarehouses(shopID int, activeOnly bool) ([]*models.Warehouse, error)
	CreateWarehouse(name, location string, shopID int) (*models.Warehouse, error)
	UpdateWarehouse(id int, name, location string, active *bool) (*models.Warehouse, error)
	TransferStock(productID, fromWarehouseID, toWarehouseID int, quantity int32) error
	GetStock(productID, warehouseID int) (*models.Stock, error)
	AddStock(productID, warehouseID int, quantity, reserved int32) (*models.Stock, error)
	SubtractStock(productID, warehouseID int, quantity, reserved int32) (*models.Stock, error)
	SetStock(productID, warehouseID int, quantity, reserved int32) (*models.Stock, error)
}
