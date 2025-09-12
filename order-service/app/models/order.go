package models

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"
)

type Order struct {
	ID          int         `gorm:"primaryKey" json:"id"`
	UserID      int         `json:"user_id"`
	Items       []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	TotalAmount float64     `json:"total_amount"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ExpiresAt   time.Time   `json:"exipres_at"`
}

type OrderItem struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	OrderID   string    `json:"order_id"`
	ProductID int       `json:"product_id"`
	ShopID    int       `json:"shop_id"`
	Quantity  int32     `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
