package service

import (
	"context"
	"errors"
	"fmt"

	"firecrest/db"
	"firecrest/internal/repository"
)

// UserService defines the interface for user business logic.
type UserService interface {
	GetUser(ctx context.Context, id int64) (db.User, error)
	CreateUser(ctx context.Context, input CreateUserInput) (db.User, error)
}

// CreateUserInput represents the input for creating a user.
type CreateUserInput struct {
	Email     string
	FirstName string
	LastName  string
	Role      db.UserRole
}

// Validate checks if the input is valid.
func (i CreateUserInput) Validate() error {
	if i.Email == "" {
		return fmt.Errorf("%w: email is required", ErrInvalidInput)
	}
	if i.FirstName == "" {
		return fmt.Errorf("%w: first_name is required", ErrInvalidInput)
	}
	if i.LastName == "" {
		return fmt.Errorf("%w: last_name is required", ErrInvalidInput)
	}
	return nil
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new UserService with the given repository.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUser(ctx context.Context, id int64) (db.User, error) {
	if id <= 0 {
		return db.User{}, fmt.Errorf("%w: invalid user id", ErrInvalidInput)
	}
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (db.User, error) {
	if err := input.Validate(); err != nil {
		return db.User{}, err
	}

	// Handle role validation - use "entrant" as the default if no role is provided
	role := input.Role
	if role == "" {
		return db.User{}, errors.New("role is required")
	}

	return s.userRepo.Create(ctx, db.CreateUserParams{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      role,
	})
}
