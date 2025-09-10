package models

import "errors"

type Error struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)
