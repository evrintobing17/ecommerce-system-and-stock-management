package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/evrintobing17/ecommerce-system/order-service/app"
	"github.com/evrintobing17/ecommerce-system/order-service/app/models"

	productProto "github.com/evrintobing17/ecommerce-system/shared/proto/product"
	warehouseProto "github.com/evrintobing17/ecommerce-system/shared/proto/warehouse"
)

// type orderUsecase struct {
// 	orderRepo app.OrderRepository
// }

type orderUsecase struct {
	orderRepo       app.OrderRepository
	productClient   productProto.ProductServiceClient
	warehouseClient warehouseProto.WarehouseServiceClient
	orderTimeout    time.Duration
}

func NewOrderUsecase(orderRepo app.OrderRepository, productClient productProto.ProductServiceClient, warehouseClient warehouseProto.WarehouseServiceClient, orderTimeout time.Duration) app.OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo,
		productClient:   productClient,
		warehouseClient: warehouseClient,
		orderTimeout:    orderTimeout,
	}
}

func (u *orderUsecase) Checkout(userID int, items []models.OrderItem) (*models.Order, error) {
	// 1. Validate products and check stock availability
	for _, item := range items {
		// Check product exists and get current stock
		product, err := u.productClient.GetProduct(context.Background(), &productProto.GetProductRequest{
			ProductId: int32(item.ProductID),
		})
		if err != nil {
			return nil, fmt.Errorf("product %s not found: %w", item.ProductID, err)
		}

		// Check if enough stock is available
		if product.Product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", item.ProductID)
		}
	}

	// 2. Reserve stock in warehouse
	for _, item := range items {
		// Find available warehouses with stock
		warehouses, err := u.warehouseClient.GetWarehouses(context.Background(), &warehouseProto.GetWarehousesRequest{
			ActiveOnly: true,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get warehouses: %w", err)
		}

		// Try to reserve stock in any available warehouse
		stockReserved := false
		for _, warehouse := range warehouses.Warehouses {
			// Check stock in this warehouse
			stock, err := u.warehouseClient.GetStock(context.Background(), &warehouseProto.GetStockRequest{
				ProductId:   int32(item.ProductID),
				WarehouseId: warehouse.Id,
			})
			if err != nil {
				continue // Skip if error checking stock
			}

			availableStock := stock.Stock.Quantity - stock.Stock.Reserved
			if availableStock >= item.Quantity {
				// Reserve the stock
				_, err := u.warehouseClient.UpdateStock(context.Background(), &warehouseProto.UpdateStockRequest{
					ProductId:   int32(item.ProductID),
					WarehouseId: warehouse.Id,
					Reserved:    item.Quantity,
					Operation:   "add_reserved",
				})
				if err == nil {
					stockReserved = true
					break
				}
			}
		}

		if !stockReserved {
			return nil, fmt.Errorf("could not reserve stock for product %s", item.ProductID)
		}
	}

	// 3. Calculate total amount
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// 4. Create order with expiration time
	expiresAt := time.Now().Add(u.orderTimeout)

	// Convert to repository items
	var repoItems []models.OrderItem
	for _, item := range items {
		repoItems = append(repoItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	order := &models.Order{
		UserID:      userID,
		Items:       repoItems,
		TotalAmount: totalAmount,
		Status:      models.OrderStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
	}

	err := u.orderRepo.Create(order)
	if err != nil {
		// If order creation fails, release reserved stock
		u.releaseReservedStock(items)
		return nil, err
	}

	// 5. Return the created order
	var resultItems []models.OrderItem
	for _, item := range order.Items {
		resultItems = append(resultItems, models.OrderItem{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &models.Order{
		ID:          order.ID,
		UserID:      order.UserID,
		Items:       resultItems,
		TotalAmount: order.TotalAmount,
		Status:      models.OrderStatus(order.Status),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		ExpiresAt:   order.ExpiresAt,
	}, nil
}

func (u *orderUsecase) releaseReservedStock(items []models.OrderItem) {
	for _, item := range items {
		// Find where stock was reserved and release it
		warehouses, err := u.warehouseClient.GetWarehouses(context.Background(), &warehouseProto.GetWarehousesRequest{
			ActiveOnly: true,
		})
		if err != nil {
			log.Printf("Error getting warehouses for stock release: %v", err)
			continue
		}

		for _, warehouse := range warehouses.Warehouses {
			// Try to release reserved stock
			_, err := u.warehouseClient.UpdateStock(context.Background(), &warehouseProto.UpdateStockRequest{
				ProductId:   int32(item.ProductID),
				WarehouseId: warehouse.Id,
				Reserved:    item.Quantity,
				Operation:   "subtract_reserved",
			})
			if err == nil {
				break // Stock released successfully
			}
		}
	}
}

func (u *orderUsecase) ProcessPayment(orderID int, paymentMethod, paymentDetails string) (*models.Order, error) {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != models.OrderStatusPending {
		return nil, errors.New("order is not in pending status")
	}

	// Process payment (simulated)
	// In a real implementation, this would integrate with a payment gateway
	order.Status = models.OrderStatusPaid
	order.UpdatedAt = time.Now()

	err = u.orderRepo.UpdateStatus(orderID, order.Status)
	if err != nil {
		return nil, err
	}

	// Convert reserved stock to actual deduction
	for _, item := range order.Items {
		// Find where stock was reserved
		warehouses, err := u.warehouseClient.GetWarehouses(context.Background(), &warehouseProto.GetWarehousesRequest{
			ActiveOnly: true,
		})
		if err != nil {
			log.Printf("Error getting warehouses for stock deduction: %v", err)
			continue
		}

		for _, warehouse := range warehouses.Warehouses {
			// Deduct the reserved stock
			_, err := u.warehouseClient.UpdateStock(context.Background(), &warehouseProto.UpdateStockRequest{
				ProductId:   int32(item.ProductID),
				WarehouseId: warehouse.Id,
				Quantity:    item.Quantity,
				Reserved:    item.Quantity,
				Operation:   "deduct_reserved",
			})
			if err == nil {
				break // Stock deducted successfully
			}
		}
	}

	// Convert to usecase order
	var items []models.OrderItem
	for _, item := range order.Items {
		items = append(items, models.OrderItem{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &models.Order{
		ID:          order.ID,
		UserID:      order.UserID,
		Items:       items,
		TotalAmount: order.TotalAmount,
		Status:      models.OrderStatus(order.Status),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		ExpiresAt:   order.ExpiresAt,
	}, nil
}

func (u *orderUsecase) ReleaseExpiredOrders() error {
	// Find orders that have expired (not paid within timeout)
	expiredTime := time.Now()
	orders, err := u.orderRepo.FindExpiredOrders(expiredTime)
	if err != nil {
		return err
	}

	for _, order := range orders {
		if order.Status == models.OrderStatusPending {
			// Release reserved stock
			for _, item := range order.Items {
				// Find where stock was reserved and release it
				warehouses, err := u.warehouseClient.GetWarehouses(context.Background(), &warehouseProto.GetWarehousesRequest{
					ActiveOnly: true,
				})
				if err != nil {
					log.Printf("Error getting warehouses for stock release: %v", err)
					continue
				}

				for _, warehouse := range warehouses.Warehouses {
					// Try to release reserved stock
					_, err := u.warehouseClient.UpdateStock(context.Background(), &warehouseProto.UpdateStockRequest{
						ProductId:   int32(item.ProductID),
						WarehouseId: warehouse.Id,
						Reserved:    item.Quantity,
						Operation:   "subtract_reserved",
					})
					if err == nil {
						break // Stock released successfully
					}
				}
			}

			// Update order status to cancelled
			order.Status = models.OrderStatusCancelled
			err := u.orderRepo.UpdateStatus(order.ID, order.Status)
			if err != nil {
				log.Printf("Error cancelling order %s: %v", order.ID, err)
			}
		}
	}

	return nil
}

func (u *orderUsecase) CreateOrder(userID int, items []models.OrderItem) (*models.Order, error) {
	// Calculate total amount
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.Price * float64(item.Quantity)
	}

	// Convert to repository items
	var repoItems []models.OrderItem
	for _, item := range items {
		repoItems = append(repoItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	order := &models.Order{
		UserID:      userID,
		Items:       repoItems,
		TotalAmount: totalAmount,
		Status:      models.OrderStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := u.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	// Convert back to usecase order
	var resultItems []models.OrderItem
	for _, item := range order.Items {
		resultItems = append(resultItems, models.OrderItem{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &models.Order{
		ID:          order.ID,
		UserID:      order.UserID,
		Items:       resultItems,
		TotalAmount: order.TotalAmount,
		Status:      models.OrderStatus(order.Status),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

func (u *orderUsecase) GetOrder(id int) (*models.Order, error) {
	order, err := u.orderRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Convert to usecase order
	var items []models.OrderItem
	for _, item := range order.Items {
		items = append(items, models.OrderItem{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}

	return &models.Order{
		ID:          order.ID,
		UserID:      order.UserID,
		Items:       items,
		TotalAmount: order.TotalAmount,
		Status:      models.OrderStatus(order.Status),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

func (u *orderUsecase) GetUserOrders(userID, page, limit int) ([]*models.Order, int64, error) {
	orders, total, err := u.orderRepo.FindByUserID(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	var result []*models.Order
	for _, order := range orders {
		var items []models.OrderItem
		for _, item := range order.Items {
			items = append(items, models.OrderItem{
				ID:        item.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Price,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
			})
		}

		result = append(result, &models.Order{
			ID:          order.ID,
			UserID:      order.UserID,
			Items:       items,
			TotalAmount: order.TotalAmount,
			Status:      models.OrderStatus(order.Status),
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		})
	}

	return result, total, nil
}


func (u *orderUsecase) CancelOrder(orderID int) error {
	order, err := u.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != models.OrderStatusPending {
		return errors.New("only pending orders can be cancelled")
	}

	order.Status = models.OrderStatusCancelled
	order.UpdatedAt = time.Now()

	return u.orderRepo.UpdateStatus(orderID, order.Status)
}
