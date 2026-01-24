package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"firecrest/db"
	"firecrest/internal/repository"
)

// mockAuthRepository implements repository.AuthRepository for testing.
type mockAuthRepository struct {
	getUserByEmailFunc          func(ctx context.Context, email string) (db.User, error)
	getCredentialsByUserIDFunc  func(ctx context.Context, userID int64) (db.AuthCredential, error)
	getCredentialsByEmailFunc   func(ctx context.Context, email string) (db.AuthCredential, error)
	isAccountLockedFunc         func(ctx context.Context, userID int64) (bool, error)
	incrementFailedAttemptsFunc func(ctx context.Context, userID int64) error
	lockAccountFunc             func(ctx context.Context, userID int64, lockUntil time.Time) error
	updateLastLoginFunc         func(ctx context.Context, userID int64) error
	verifyEmailFunc             func(ctx context.Context, userID int64) error
	createCredentialsFunc       func(ctx context.Context, userID int64, passwordHash string) (db.AuthCredential, error)
}

func (m *mockAuthRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(ctx, email)
	}
	return db.User{}, nil
}

func (m *mockAuthRepository) GetCredentialsByUserID(ctx context.Context, userID int64) (db.AuthCredential, error) {
	if m.getCredentialsByUserIDFunc != nil {
		return m.getCredentialsByUserIDFunc(ctx, userID)
	}
	return db.AuthCredential{}, nil
}

func (m *mockAuthRepository) GetCredentialsByEmail(ctx context.Context, email string) (db.AuthCredential, error) {
	if m.getCredentialsByEmailFunc != nil {
		return m.getCredentialsByEmailFunc(ctx, email)
	}
	return db.AuthCredential{}, nil
}

func (m *mockAuthRepository) IsAccountLocked(ctx context.Context, userID int64) (bool, error) {
	if m.isAccountLockedFunc != nil {
		return m.isAccountLockedFunc(ctx, userID)
	}
	return false, nil
}

func (m *mockAuthRepository) IncrementFailedAttempts(ctx context.Context, userID int64) error {
	if m.incrementFailedAttemptsFunc != nil {
		return m.incrementFailedAttemptsFunc(ctx, userID)
	}
	return nil
}

func (m *mockAuthRepository) LockAccount(ctx context.Context, userID int64, lockUntil time.Time) error {
	if m.lockAccountFunc != nil {
		return m.lockAccountFunc(ctx, userID, lockUntil)
	}
	return nil
}

func (m *mockAuthRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	if m.updateLastLoginFunc != nil {
		return m.updateLastLoginFunc(ctx, userID)
	}
	return nil
}

func (m *mockAuthRepository) VerifyEmail(ctx context.Context, userID int64) error {
	if m.verifyEmailFunc != nil {
		return m.verifyEmailFunc(ctx, userID)
	}
	return nil
}

func (m *mockAuthRepository) CreateCredentials(ctx context.Context, userID int64, passwordHash string) (db.AuthCredential, error) {
	if m.createCredentialsFunc != nil {
		return m.createCredentialsFunc(ctx, userID, passwordHash)
	}
	return db.AuthCredential{}, nil
}

// MockClock implements Clock for testing.
type MockClock struct {
	CurrentTime time.Time
}

func (m *MockClock) Now() time.Time {
	return m.CurrentTime
}

// MockHasher implements PasswordHasher for testing.
type MockHasher struct {
	CompareFunc  func(hashedPassword, password []byte) error
	GenerateFunc func(password []byte, cost int) ([]byte, error)
}

func (m *MockHasher) CompareHashAndPassword(hashedPassword, password []byte) error {
	if m.CompareFunc != nil {
		return m.CompareFunc(hashedPassword, password)
	}
	// Default: simple comparison for tests
	if string(hashedPassword) == string(password) {
		return nil
	}
	return bcrypt.ErrMismatchedHashAndPassword
}

func (m *MockHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(password, cost)
	}
	// Default: return password as-is for tests
	return password, nil
}

func TestAuthService_SignIn(t *testing.T) {
	t.Run("returns auth result for valid credentials", func(t *testing.T) {
		expectedUser := db.User{
			ID:        1,
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
		}

		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return expectedUser, nil
			},
			getCredentialsByUserIDFunc: func(ctx context.Context, userID int64) (db.AuthCredential, error) {
				return db.AuthCredential{
					UserID:       1,
					PasswordHash: "hashed_password",
					EmailVerifiedAt: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
					FailedLoginAttempts: 0,
				}, nil
			},
			isAccountLockedFunc: func(ctx context.Context, userID int64) (bool, error) {
				return false, nil
			},
			updateLastLoginFunc: func(ctx context.Context, userID int64) error {
				return nil
			},
		}

		userRepo := &mockUserRepository{}

		hasher := &MockHasher{
			CompareFunc: func(hashedPassword, password []byte) error {
				if string(hashedPassword) == "hashed_password" && string(password) == "correct_password" {
					return nil
				}
				return bcrypt.ErrMismatchedHashAndPassword
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: userRepo,
			clock:    RealClock{},
			hasher:   hasher,
		}

		result, err := svc.SignIn(context.Background(), SignInInput{
			Email:      "test@example.com",
			Password:   "correct_password",
			RememberMe: true,
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.User.ID != expectedUser.ID {
			t.Errorf("expected user ID %d, got %d", expectedUser.ID, result.User.ID)
		}
		if !result.RememberMe {
			t.Error("expected RememberMe to be true")
		}
	})

	t.Run("normalizes email to lowercase and trimmed", func(t *testing.T) {
		var capturedEmail string
		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				capturedEmail = email
				return db.User{ID: 1, Email: email}, nil
			},
			getCredentialsByUserIDFunc: func(ctx context.Context, userID int64) (db.AuthCredential, error) {
				return db.AuthCredential{
					UserID:       1,
					PasswordHash: "hashed_password",
					EmailVerifiedAt: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
				}, nil
			},
			isAccountLockedFunc: func(ctx context.Context, userID int64) (bool, error) {
				return false, nil
			},
		}

		hasher := &MockHasher{
			CompareFunc: func(hashedPassword, password []byte) error {
				return nil
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   hasher,
		}

		_, _ = svc.SignIn(context.Background(), SignInInput{
			Email:    "  Test@EXAMPLE.com  ",
			Password: "password",
		})

		if capturedEmail != "test@example.com" {
			t.Errorf("expected normalized email 'test@example.com', got '%s'", capturedEmail)
		}
	})

	t.Run("returns ErrInvalidInput for empty email", func(t *testing.T) {
		svc := &authService{
			authRepo: &mockAuthRepository{},
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   &MockHasher{},
		}

		_, err := svc.SignIn(context.Background(), SignInInput{
			Email:    "",
			Password: "password",
		})

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrInvalidInput for empty password", func(t *testing.T) {
		svc := &authService{
			authRepo: &mockAuthRepository{},
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   &MockHasher{},
		}

		_, err := svc.SignIn(context.Background(), SignInInput{
			Email:    "test@example.com",
			Password: "",
		})

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrInvalidCredentials for non-existent user", func(t *testing.T) {
		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return db.User{}, repository.ErrNotFound
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   &MockHasher{},
		}

		_, err := svc.SignIn(context.Background(), SignInInput{
			Email:    "nonexistent@example.com",
			Password: "password",
		})

		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("returns ErrInvalidCredentials for wrong password", func(t *testing.T) {
		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return db.User{ID: 1, Email: email}, nil
			},
			getCredentialsByUserIDFunc: func(ctx context.Context, userID int64) (db.AuthCredential, error) {
				return db.AuthCredential{
					UserID:              1,
					PasswordHash:        "hashed_password",
					FailedLoginAttempts: 0,
				}, nil
			},
			isAccountLockedFunc: func(ctx context.Context, userID int64) (bool, error) {
				return false, nil
			},
			incrementFailedAttemptsFunc: func(ctx context.Context, userID int64) error {
				return nil
			},
		}

		hasher := &MockHasher{
			CompareFunc: func(hashedPassword, password []byte) error {
				return bcrypt.ErrMismatchedHashAndPassword
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   hasher,
		}

		_, err := svc.SignIn(context.Background(), SignInInput{
			Email:    "test@example.com",
			Password: "wrong_password",
		})

		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("increments failed attempts on wrong password", func(t *testing.T) {
		var incrementCalled bool
		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return db.User{ID: 1}, nil
			},
			getCredentialsByUserIDFunc: func(ctx context.Context, userID int64) (db.AuthCredential, error) {
				return db.AuthCredential{
					UserID:              1,
					PasswordHash:        "hashed_password",
					FailedLoginAttempts: 2,
				}, nil
			},
			isAccountLockedFunc: func(ctx context.Context, userID int64) (bool, error) {
				return false, nil
			},
			incrementFailedAttemptsFunc: func(ctx context.Context, userID int64) error {
				incrementCalled = true
				return nil
			},
		}

		hasher := &MockHasher{
			CompareFunc: func(hashedPassword, password []byte) error {
				return bcrypt.ErrMismatchedHashAndPassword
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   hasher,
		}

		_, _ = svc.SignIn(context.Background(), SignInInput{
			Email:    "test@example.com",
			Password: "wrong_password",
		})

		if !incrementCalled {
			t.Error("expected IncrementFailedAttempts to be called")
		}
	})

	t.Run("locks account on 5th failed attempt with correct lockUntil time", func(t *testing.T) {
		mockTime := time.Date(2026, 1, 23, 12, 0, 0, 0, time.UTC)
		expectedLockUntil := mockTime.Add(AccountLockoutDuration)
		var capturedLockUntil time.Time

		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return db.User{ID: 1}, nil
			},
			getCredentialsByUserIDFunc: func(ctx context.Context, userID int64) (db.AuthCredential, error) {
				return db.AuthCredential{
					UserID:              1,
					PasswordHash:        "hashed_password",
					FailedLoginAttempts: 4, // This will be the 5th attempt
				}, nil
			},
			isAccountLockedFunc: func(ctx context.Context, userID int64) (bool, error) {
				return false, nil
			},
			incrementFailedAttemptsFunc: func(ctx context.Context, userID int64) error {
				return nil
			},
			lockAccountFunc: func(ctx context.Context, userID int64, lockUntil time.Time) error {
				capturedLockUntil = lockUntil
				return nil
			},
		}

		hasher := &MockHasher{
			CompareFunc: func(hashedPassword, password []byte) error {
				return bcrypt.ErrMismatchedHashAndPassword
			},
		}

		clock := &MockClock{CurrentTime: mockTime}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    clock,
			hasher:   hasher,
		}

		_, err := svc.SignIn(context.Background(), SignInInput{
			Email:    "test@example.com",
			Password: "wrong_password",
		})

		if !errors.Is(err, ErrAccountLocked) {
			t.Errorf("expected ErrAccountLocked, got %v", err)
		}

		if !capturedLockUntil.Equal(expectedLockUntil) {
			t.Errorf("expected lock until %v, got %v", expectedLockUntil, capturedLockUntil)
		}
	})

	t.Run("returns ErrAccountLocked for locked account", func(t *testing.T) {
		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return db.User{ID: 1}, nil
			},
			getCredentialsByUserIDFunc: func(ctx context.Context, userID int64) (db.AuthCredential, error) {
				return db.AuthCredential{
					UserID:              1,
					PasswordHash:        "hashed_password",
					FailedLoginAttempts: 5,
				}, nil
			},
			isAccountLockedFunc: func(ctx context.Context, userID int64) (bool, error) {
				return true, nil
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   &MockHasher{},
		}

		_, err := svc.SignIn(context.Background(), SignInInput{
			Email:    "test@example.com",
			Password: "password",
		})

		if !errors.Is(err, ErrAccountLocked) {
			t.Errorf("expected ErrAccountLocked, got %v", err)
		}
	})

	t.Run("returns ErrEmailNotVerified for unverified email", func(t *testing.T) {
		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return db.User{ID: 1}, nil
			},
			getCredentialsByUserIDFunc: func(ctx context.Context, userID int64) (db.AuthCredential, error) {
				return db.AuthCredential{
					UserID:       1,
					PasswordHash: "hashed_password",
					EmailVerifiedAt: pgtype.Timestamptz{
						Valid: false, // Not verified
					},
					FailedLoginAttempts: 0,
				}, nil
			},
			isAccountLockedFunc: func(ctx context.Context, userID int64) (bool, error) {
				return false, nil
			},
		}

		hasher := &MockHasher{
			CompareFunc: func(hashedPassword, password []byte) error {
				return nil // Password is correct
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   hasher,
		}

		_, err := svc.SignIn(context.Background(), SignInInput{
			Email:    "test@example.com",
			Password: "password",
		})

		if !errors.Is(err, ErrEmailNotVerified) {
			t.Errorf("expected ErrEmailNotVerified, got %v", err)
		}
	})

	t.Run("calls UpdateLastLogin on successful sign in", func(t *testing.T) {
		var updateLastLoginCalled bool
		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return db.User{ID: 1}, nil
			},
			getCredentialsByUserIDFunc: func(ctx context.Context, userID int64) (db.AuthCredential, error) {
				return db.AuthCredential{
					UserID:       1,
					PasswordHash: "hashed_password",
					EmailVerifiedAt: pgtype.Timestamptz{
						Time:  time.Now(),
						Valid: true,
					},
				}, nil
			},
			isAccountLockedFunc: func(ctx context.Context, userID int64) (bool, error) {
				return false, nil
			},
			updateLastLoginFunc: func(ctx context.Context, userID int64) error {
				updateLastLoginCalled = true
				return nil
			},
		}

		hasher := &MockHasher{
			CompareFunc: func(hashedPassword, password []byte) error {
				return nil
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   hasher,
		}

		_, _ = svc.SignIn(context.Background(), SignInInput{
			Email:    "test@example.com",
			Password: "password",
		})

		if !updateLastLoginCalled {
			t.Error("expected UpdateLastLogin to be called")
		}
	})

	t.Run("propagates repository errors from GetUserByEmail", func(t *testing.T) {
		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return db.User{}, errors.New("database error")
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   &MockHasher{},
		}

		_, err := svc.SignIn(context.Background(), SignInInput{
			Email:    "test@example.com",
			Password: "password",
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if errors.Is(err, ErrInvalidCredentials) {
			t.Error("should not wrap database errors as invalid credentials")
		}
	})

	t.Run("propagates repository errors from GetCredentialsByUserID", func(t *testing.T) {
		authRepo := &mockAuthRepository{
			getUserByEmailFunc: func(ctx context.Context, email string) (db.User, error) {
				return db.User{ID: 1}, nil
			},
			getCredentialsByUserIDFunc: func(ctx context.Context, userID int64) (db.AuthCredential, error) {
				return db.AuthCredential{}, errors.New("database error")
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   &MockHasher{},
		}

		_, err := svc.SignIn(context.Background(), SignInInput{
			Email:    "test@example.com",
			Password: "password",
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if errors.Is(err, ErrInvalidCredentials) {
			t.Error("should not wrap database errors as invalid credentials")
		}
	})
}

func TestAuthService_VerifyEmail(t *testing.T) {
	t.Run("succeeds for valid user ID", func(t *testing.T) {
		authRepo := &mockAuthRepository{
			verifyEmailFunc: func(ctx context.Context, userID int64) error {
				return nil
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   &MockHasher{},
		}

		err := svc.VerifyEmail(context.Background(), 1)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("propagates repository errors", func(t *testing.T) {
		authRepo := &mockAuthRepository{
			verifyEmailFunc: func(ctx context.Context, userID int64) error {
				return errors.New("database error")
			},
		}

		svc := &authService{
			authRepo: authRepo,
			userRepo: &mockUserRepository{},
			clock:    RealClock{},
			hasher:   &MockHasher{},
		}

		err := svc.VerifyEmail(context.Background(), 1)

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
