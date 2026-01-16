package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"firecrest/db"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (db.User, error)
	Create(ctx context.Context, params db.CreateUserParams) (db.User, error)
}

type userRepository struct {
	queries *db.Queries
}

// NewUserRepository creates a new UserRepository backed by the given queries.
func NewUserRepository(queries *db.Queries) UserRepository {
	return &userRepository{queries: queries}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (db.User, error) {
	user, err := r.queries.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, params)
}
