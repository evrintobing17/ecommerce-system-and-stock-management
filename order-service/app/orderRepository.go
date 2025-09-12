package app

import (
	"time"

	"github.com/evrintobing17/ecommerce-system/order-service/app/models"
)

type OrderRepository interface {
	FindExpiredOrders(expiredTime time.Time) ([]*models.Order, error) // New method
	Create(order *models.Order) error
	FindByID(id int) (*models.Order, error)
	FindByUserID(userID, page, limit int) ([]*models.Order, int64, error)
	Update(order *models.Order) error
	UpdateStatus(id int, status models.OrderStatus) error
	Delete(id int) error
}
