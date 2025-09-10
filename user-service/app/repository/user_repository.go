package repository

import (
	"context"
	"database/sql"
	user "user-service/app"
	"user-service/app/models"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, phone, password_hash 
		FROM users 
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) FindByPhoneOrEmail(ctx context.Context, data string) (*models.User, error) {
	query := `
		SELECT id, email, phone, password_hash 
		FROM users 
		WHERE email = $1 OR phone = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, data).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
