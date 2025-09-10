package app

import (
	"context"
	"user-service/app/models"
)

type UserUsecase interface {
	Login(ctx context.Context, data, password string) (string, *models.UserResponse, error)
}
