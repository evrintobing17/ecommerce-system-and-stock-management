package repository

import (
	"errors"
	"time"

	"github.com/evrintobing17/ecommerce-system/order-service/app"
	"github.com/evrintobing17/ecommerce-system/order-service/app/models"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) app.OrderRepository {
	return &orderRepository{db: db}
}


func (r *orderRepository) FindExpiredOrders(expiredTime time.Time) ([]*models.Order, error) {
	var orders []*models.Order
	err := r.db.Preload("Items").
		Where("status = ? AND expires_at <= ?", models.OrderStatusPending, expiredTime).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id int) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("Items").First(&order, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID, page, limit int) ([]*models.Order, int64, error) {
	var orders []*models.Order
	var total int64

	// Get total count
	err := r.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err = r.db.Preload("Items").Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) UpdateStatus(id int, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) Delete(id int) error {
	return r.db.Delete(&models.Order{}, "id = ?", id).Error
}
