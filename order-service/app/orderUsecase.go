package app

import "github.com/evrintobing17/ecommerce-system/order-service/app/models"

type OrderUsecase interface {
	CreateOrder(userID int, items []models.OrderItem) (*models.Order, error)
	GetOrder(id int) (*models.Order, error)
	GetUserOrders(userID, page, limit int) ([]*models.Order, int64, error)
	ProcessPayment(orderID int, paymentMethod, paymentDetails string) (*models.Order, error)
	CancelOrder(orderID int) error
	Checkout(userID int, items []models.OrderItem) (*models.Order, error)
	ReleaseExpiredOrders() error
}
