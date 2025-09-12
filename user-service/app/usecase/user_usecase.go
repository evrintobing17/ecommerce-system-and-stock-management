package usecase

import (
	"errors"
	"time"

	"github.com/evrintobing17/ecommerce-system/shared"
	"github.com/evrintobing17/ecommerce-system/user-service/app"
	"github.com/evrintobing17/ecommerce-system/user-service/app/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo  app.UserRepository
	jwtSecret string
}

func NewUserUsecase(userRepo app.UserRepository, jwtSecret string) app.UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) Register(email, phone, password, name string) (*models.User, string, error) {
	// Check if user already exists
	_, err := u.userRepo.FindByEmail(email)
	if err == nil {
		return nil, "", errors.New("user with this email already exists")
	}

	_, err = u.userRepo.FindByPhone(phone)
	if err == nil {
		return nil, "", errors.New("user with this phone already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := &models.User{
		Email:     email,
		Phone:     phone,
		Name:      name,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.userRepo.Create(user)
	if err != nil {
		return nil, "", err
	}

	// Generate JWT token
	token, err := shared.GenerateToken(user.ID, user.Email, u.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return &models.User{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, token, nil
}

func (u *userUsecase) Login(emailOrPhone, password string) (*models.User, string, error) {
	var user *models.User
	var err error

	// Try to find by email first
	user, err = u.userRepo.FindByEmail(emailOrPhone)
	if err != nil {
		// If not found by email, try by phone
		user, err = u.userRepo.FindByPhone(emailOrPhone)
		if err != nil {
			return nil, "", errors.New("invalid credentials")
		}
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := shared.GenerateToken(user.ID, user.Email, u.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return &models.User{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, token, nil
}

func (u *userUsecase) ValidateToken(token string) (bool, *models.User, error) {
	claims, err := u.parseToken(token)
	if err != nil {
		return false, nil, err
	}

	user, err := u.userRepo.FindByID(claims.UserID)
	if err != nil {
		return false, nil, errors.New("user not found")
	}

	return true, &models.User{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (u *userUsecase) GetUser(id int) (*models.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        user.ID,
		Email:     user.Email,
		Phone:     user.Phone,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (u *userUsecase) UpdateUser(user *models.User) error {
	existingUser, err := u.userRepo.FindByID(user.ID)
	if err != nil {
		return err
	}

	existingUser.Name = user.Name
	existingUser.Email = user.Email
	existingUser.Phone = user.Phone
	existingUser.UpdatedAt = time.Now()

	return u.userRepo.Update(existingUser)
}

func (u *userUsecase) parseToken(tokenString string) (*shared.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &shared.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*shared.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
