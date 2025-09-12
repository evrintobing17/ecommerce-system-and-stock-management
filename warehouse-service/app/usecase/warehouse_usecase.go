package usecase

import (
	"context"
	"errors"
	"time"

	protoShop "github.com/evrintobing17/ecommerce-system/shared/proto/shop"
	"github.com/evrintobing17/ecommerce-system/warehouse-service/app"
	"github.com/evrintobing17/ecommerce-system/warehouse-service/app/models"
)

type warehouseUsecase struct {
	warehouseRepo app.WarehouseRepository
	stockRepo     app.StockRepository
	shopProto     protoShop.ShopServiceClient
}

func NewWarehouseUsecase(warehouseRepo app.WarehouseRepository, stockRepo app.StockRepository, shopProto protoShop.ShopServiceClient) app.WarehouseUsecase {
	return &warehouseUsecase{
		warehouseRepo: warehouseRepo,
		stockRepo:     stockRepo,
		shopProto:     shopProto,
	}
}

func (u *warehouseUsecase) GetWarehouse(id int) (*models.Warehouse, error) {
	warehouse, err := u.warehouseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &models.Warehouse{
		ID:        warehouse.ID,
		Name:      warehouse.Name,
		Location:  warehouse.Location,
		ShopID:    warehouse.ShopID,
		Active:    warehouse.Active,
		CreatedAt: warehouse.CreatedAt,
		UpdatedAt: warehouse.UpdatedAt,
	}, nil
}

func (u *warehouseUsecase) GetWarehouses(shopID int, activeOnly bool) ([]*models.Warehouse, error) {
	warehouses, err := u.warehouseRepo.FindByShopID(shopID, activeOnly)
	if err != nil {
		return nil, err
	}

	var result []*models.Warehouse
	for _, warehouse := range warehouses {
		result = append(result, &models.Warehouse{
			ID:        warehouse.ID,
			Name:      warehouse.Name,
			Location:  warehouse.Location,
			ShopID:    warehouse.ShopID,
			Active:    warehouse.Active,
			CreatedAt: warehouse.CreatedAt,
			UpdatedAt: warehouse.UpdatedAt,
		})
	}

	return result, nil
}

func (u *warehouseUsecase) CreateWarehouse(name, location string, shopID int) (*models.Warehouse, error) {
	warehouse := &models.Warehouse{
		Name:      name,
		Location:  location,
		ShopID:    shopID,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := u.warehouseRepo.Create(warehouse)
	if err != nil {
		return nil, err
	}

	_, err = u.shopProto.GetShop(context.Background(), &protoShop.GetShopRequest{ShopId: int32(shopID)})
	if err != nil {
		return nil, err
	}

	return &models.Warehouse{
		ID:        warehouse.ID,
		Name:      warehouse.Name,
		Location:  warehouse.Location,
		ShopID:    warehouse.ShopID,
		Active:    warehouse.Active,
		CreatedAt: warehouse.CreatedAt,
		UpdatedAt: warehouse.UpdatedAt,
	}, nil
}

func (u *warehouseUsecase) UpdateWarehouse(id int, name, location string, active *bool) (*models.Warehouse, error) {
	warehouse, err := u.warehouseRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		warehouse.Name = name
	}
	if location != "" {
		warehouse.Location = location
	}
	if active != nil {
		warehouse.Active = *active
	}
	warehouse.UpdatedAt = time.Now()

	err = u.warehouseRepo.Update(warehouse)
	if err != nil {
		return nil, err
	}

	return &models.Warehouse{
		ID:        warehouse.ID,
		Name:      warehouse.Name,
		Location:  warehouse.Location,
		ShopID:    warehouse.ShopID,
		Active:    warehouse.Active,
		CreatedAt: warehouse.CreatedAt,
		UpdatedAt: warehouse.UpdatedAt,
	}, nil
}

func (u *warehouseUsecase) TransferStock(productID, fromWarehouseID, toWarehouseID int, quantity int32) error {
	// Check if from warehouse exists and is active
	fromWarehouse, err := u.warehouseRepo.FindByID(fromWarehouseID)
	if err != nil {
		return err
	}
	if !fromWarehouse.Active {
		return errors.New("source warehouse is not active")
	}

	// Check if to warehouse exists and is active
	toWarehouse, err := u.warehouseRepo.FindByID(toWarehouseID)
	if err != nil {
		return err
	}
	if !toWarehouse.Active {
		return errors.New("destination warehouse is not active")
	}

	// Check if there's enough stock in the source warehouse
	stock, err := u.stockRepo.Find(productID, fromWarehouseID)
	if err != nil {
		return err
	}

	if stock.Quantity-stock.Reserved < quantity {
		return errors.New("insufficient stock in source warehouse")
	}

	// Perform the transfer
	return u.stockRepo.Transfer(productID, fromWarehouseID, toWarehouseID, quantity)
}

func (u *warehouseUsecase) GetStock(productID, warehouseID int) (*models.Stock, error) {
	stock, err := u.stockRepo.Find(productID, warehouseID)
	if err != nil {
		return nil, err
	}

	return &models.Stock{
		ProductID:   stock.ProductID,
		WarehouseID: stock.WarehouseID,
		Quantity:    stock.Quantity,
		Reserved:    stock.Reserved,
		CreatedAt:   stock.CreatedAt,
		UpdatedAt:   stock.UpdatedAt,
	}, nil
}

func (u *warehouseUsecase) AddStock(productID, warehouseID int, quantity, reserved int32) (*models.Stock, error) {
	stock, err := u.stockRepo.Find(productID, warehouseID)
	if err != nil {
		// Stock doesn't exist, create it
		stock = &models.Stock{
			ProductID:   productID,
			WarehouseID: warehouseID,
			Quantity:    quantity,
			Reserved:    reserved,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = u.stockRepo.Create(stock)
		if err != nil {
			return nil, err
		}
	} else {
		// Stock exists, update it
		stock.Quantity += quantity
		stock.Reserved += reserved
		stock.UpdatedAt = time.Now()
		err = u.stockRepo.Update(stock)
		if err != nil {
			return nil, err
		}
	}

	return &models.Stock{
		ProductID:   stock.ProductID,
		WarehouseID: stock.WarehouseID,
		Quantity:    stock.Quantity,
		Reserved:    stock.Reserved,
		CreatedAt:   stock.CreatedAt,
		UpdatedAt:   stock.UpdatedAt,
	}, nil
}

func (u *warehouseUsecase) SubtractStock(productID, warehouseID int, quantity, reserved int32) (*models.Stock, error) {
	stock, err := u.stockRepo.Find(productID, warehouseID)
	if err != nil {
		return nil, errors.New("stock not found")
	}

	if stock.Quantity < quantity {
		return nil, errors.New("insufficient quantity")
	}

	if stock.Reserved < reserved {
		return nil, errors.New("insufficient reserved stock")
	}

	stock.Quantity -= quantity
	stock.Reserved -= reserved
	stock.UpdatedAt = time.Now()

	err = u.stockRepo.Update(stock)
	if err != nil {
		return nil, err
	}

	return &models.Stock{
		ProductID:   stock.ProductID,
		WarehouseID: stock.WarehouseID,
		Quantity:    stock.Quantity,
		Reserved:    stock.Reserved,
		CreatedAt:   stock.CreatedAt,
		UpdatedAt:   stock.UpdatedAt,
	}, nil
}

func (u *warehouseUsecase) SetStock(productID, warehouseID int, quantity, reserved int32) (*models.Stock, error) {
	stock, err := u.stockRepo.Find(productID, warehouseID)
	if err != nil {
		// Stock doesn't exist, create it
		stock = &models.Stock{
			ProductID:   productID,
			WarehouseID: warehouseID,
			Quantity:    quantity,
			Reserved:    reserved,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		err = u.stockRepo.Create(stock)
		if err != nil {
			return nil, err
		}
	} else {
		// Stock exists, update it
		stock.Quantity = quantity
		stock.Reserved = reserved
		stock.UpdatedAt = time.Now()
		err = u.stockRepo.Update(stock)
		if err != nil {
			return nil, err
		}
	}

	return &models.Stock{
		ProductID:   stock.ProductID,
		WarehouseID: stock.WarehouseID,
		Quantity:    stock.Quantity,
		Reserved:    stock.Reserved,
		CreatedAt:   stock.CreatedAt,
		UpdatedAt:   stock.UpdatedAt,
	}, nil
}
