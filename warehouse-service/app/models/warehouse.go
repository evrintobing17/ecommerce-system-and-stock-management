package models

import "time"

type Warehouse struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	ShopID    int       `json:"shop_id"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Stock struct {
	ProductID   int       `gorm:"primaryKey" json:"product_id"`
	WarehouseID int       `gorm:"primaryKey" json:"warehouse_id"`
	Quantity    int32     `json:"quantity"`
	Reserved    int32     `json:"reserved"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
