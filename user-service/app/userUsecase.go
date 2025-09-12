package app

import "github.com/evrintobing17/ecommerce-system/user-service/app/models"



type UserUsecase interface {
	Register(email, phone, password, name string) (*models.User, string, error)
	Login(emailOrPhone, password string) (*models.User, string, error)
	ValidateToken(token string) (bool, *models.User, error)
	GetUser(id int) (*models.User, error)
	UpdateUser(user *models.User) error
}
