package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"omitempty,required_without=Phone,email"`
	Phone    string `json:"phone" binding:"omitempty,required_without=Email,number"`
	Password string `json:"password" binding:"required"`
}
