package usecase

import (
	"context"
	"shared"
	"time"

	user "user-service/app"
	"user-service/app/models"

	"github.com/golang-jwt/jwt/v4"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo  user.UserRepository
	jwtSecret string
	log       shared.Log
}

func NewUserUsecase(userRepo user.UserRepository, jwtSecret string, log shared.Log) user.UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		log:       log,
	}
}

func (s *userUsecase) Login(ctx context.Context, email, password string) (string, *models.UserResponse, error) {
	user, err := s.userRepo.FindByPhoneOrEmail(ctx, email)
	if err != nil {
		s.log.ErrorLog(err)
		return "", nil, models.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		s.log.ErrorLog(models.ErrInvalidCredentials)
		return "", nil, models.ErrInvalidCredentials
	}

	token, err := s.generateToken(user)
	if err != nil {
		s.log.ErrorLog(err)
		return "", nil, err
	}

	result := &models.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Phone: user.Phone,
	}

	return token, result, nil
}

func (s *userUsecase) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
