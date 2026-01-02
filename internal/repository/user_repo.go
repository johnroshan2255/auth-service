package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/johnroshan2255/auth-service/internal/model"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
}

type PostgresUserRepo struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepo(db *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	row := r.db.QueryRow(ctx,
		"SELECT id, email, password_hash, tenant_id, role, created_at, updated_at FROM users WHERE email=$1", email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.TenantID, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r *PostgresUserRepo) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO users (id, email, password_hash, tenant_id, role, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)",
		user.ID, user.Email, user.PasswordHash, user.TenantID, user.Role, time.Now(), time.Now())
	return err
}
