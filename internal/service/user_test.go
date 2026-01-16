package service

import (
	"context"
	"errors"
	"testing"

	"firecrest/db"
	"firecrest/internal/repository"
)

// mockUserRepository implements repository.UserRepository for testing.
type mockUserRepository struct {
	getByIDFunc func(ctx context.Context, id int64) (db.User, error)
	createFunc  func(ctx context.Context, params db.CreateUserParams) (db.User, error)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id int64) (db.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return db.User{}, nil
}

func (m *mockUserRepository) Create(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, params)
	}
	return db.User{}, nil
}

func TestUserService_GetUser(t *testing.T) {
	t.Run("returns user for valid id", func(t *testing.T) {
		expected := db.User{ID: 1, Email: "test@example.com", FirstName: "Test", LastName: "User"}

		repo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id int64) (db.User, error) {
				if id == 1 {
					return expected, nil
				}
				return db.User{}, repository.ErrNotFound
			},
		}

		svc := NewUserService(repo)
		user, err := svc.GetUser(context.Background(), 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != expected.ID {
			t.Errorf("expected user ID %d, got %d", expected.ID, user.ID)
		}
	})

	t.Run("returns ErrInvalidInput for invalid id", func(t *testing.T) {
		repo := &mockUserRepository{}
		svc := NewUserService(repo)

		_, err := svc.GetUser(context.Background(), 0)

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrInvalidInput for negative id", func(t *testing.T) {
		repo := &mockUserRepository{}
		svc := NewUserService(repo)

		_, err := svc.GetUser(context.Background(), -1)

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrNotFound for non-existent user", func(t *testing.T) {
		repo := &mockUserRepository{
			getByIDFunc: func(ctx context.Context, id int64) (db.User, error) {
				return db.User{}, repository.ErrNotFound
			},
		}

		svc := NewUserService(repo)
		_, err := svc.GetUser(context.Background(), 999)

		if !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestUserService_CreateUser(t *testing.T) {
	t.Run("creates user with valid input", func(t *testing.T) {
		expected := db.User{ID: 1, Email: "test@example.com", FirstName: "Test", LastName: "User", Role: "entrant"}

		repo := &mockUserRepository{
			createFunc: func(ctx context.Context, params db.CreateUserParams) (db.User, error) {
				return expected, nil
			},
		}

		svc := NewUserService(repo)
		user, err := svc.CreateUser(context.Background(), CreateUserInput{
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Role:      "entrant",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.ID != expected.ID {
			t.Errorf("expected user ID %d, got %d", expected.ID, user.ID)
		}
	})

	t.Run("returns ErrInvalidInput for missing email", func(t *testing.T) {
		repo := &mockUserRepository{}
		svc := NewUserService(repo)

		_, err := svc.CreateUser(context.Background(), CreateUserInput{
			Email:     "",
			FirstName: "Test",
			LastName:  "User",
			Role:      "entrant",
		})

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrInvalidInput for missing first_name", func(t *testing.T) {
		repo := &mockUserRepository{}
		svc := NewUserService(repo)

		_, err := svc.CreateUser(context.Background(), CreateUserInput{
			Email:     "test@example.com",
			FirstName: "",
			LastName:  "User",
			Role:      "entrant",
		})

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrInvalidInput for missing last_name", func(t *testing.T) {
		repo := &mockUserRepository{}
		svc := NewUserService(repo)

		_, err := svc.CreateUser(context.Background(), CreateUserInput{
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "",
			Role:      "entrant",
		})

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns error for missing role", func(t *testing.T) {
		repo := &mockUserRepository{}
		svc := NewUserService(repo)

		_, err := svc.CreateUser(context.Background(), CreateUserInput{
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Role:      "",
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("propagates repository errors", func(t *testing.T) {
		repo := &mockUserRepository{
			createFunc: func(ctx context.Context, params db.CreateUserParams) (db.User, error) {
				return db.User{}, errors.New("database error")
			},
		}

		svc := NewUserService(repo)
		_, err := svc.CreateUser(context.Background(), CreateUserInput{
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Role:      "entrant",
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
