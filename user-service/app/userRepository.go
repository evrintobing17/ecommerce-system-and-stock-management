package app

import (
	"context"
	"user-service/app/models"
)

type UserRepository interface {
	FindByID(ctx context.Context, id int) (*models.User, error)
	FindByPhoneOrEmail(ctx context.Context, input string) (*models.User, error)
}
