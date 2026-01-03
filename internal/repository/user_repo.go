package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"github.com/johnroshan2255/auth-service/internal/model"
)

type UserRepository interface {
	GetByID(ctx context.Context, userUUID string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type PostgresUserRepo struct {
	db *gorm.DB
}

func NewPostgresUserRepo(db *gorm.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, userUUID string) (*model.User, error) {
	user := &model.User{}
	err := r.db.WithContext(ctx).Where("uuid = ?", userUUID).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := r.db.WithContext(ctx).Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.WithContext(ctx).Where("username = ?", username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *PostgresUserRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if email already exists
		var emailCount int64
		if err := tx.Model(&model.User{}).Where("email = ?", user.Email).Count(&emailCount).Error; err != nil {
			return fmt.Errorf("failed to check email: %w", err)
		}
		if emailCount > 0 {
			return errors.New("email already exists")
		}

		// Check if username already exists
		var usernameCount int64
		if err := tx.Model(&model.User{}).Where("username = ?", user.Username).Count(&usernameCount).Error; err != nil {
			return fmt.Errorf("failed to check username: %w", err)
		}
		if usernameCount > 0 {
			return errors.New("username already exists")
		}

		// Create user (UUID and timestamps are handled by GORM hooks/defaults)
		if err := tx.Create(user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		return nil
	})
}
