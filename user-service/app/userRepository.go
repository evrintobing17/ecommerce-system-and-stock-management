package app

import "github.com/evrintobing17/ecommerce-system/user-service/app/models"



type UserRepository interface {
	Create(user *models.User) error
	FindByID(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
}
