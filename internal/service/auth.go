package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"firecrest/db"
	"firecrest/internal/repository"
)

// Authentication errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailNotVerified   = errors.New("email address not verified")
	ErrAccountLocked      = errors.New("account is locked due to too many failed login attempts")
	ErrEmailExists        = errors.New("email address already registered")
)

// Authentication constants
const (
	BcryptCost             = 12
	MaxLoginAttempts       = 5
	AccountLockoutDuration = 15 * time.Minute
	MinPasswordLength      = 8
)

// Clock provides time-related operations for testing.
type Clock interface {
	Now() time.Time
}

// RealClock implements Clock using the standard time package.
type RealClock struct{}

// Now returns the current time.
func (RealClock) Now() time.Time {
	return time.Now()
}

// PasswordHasher provides password hashing operations for testing.
type PasswordHasher interface {
	CompareHashAndPassword(hashedPassword, password []byte) error
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
}

// BcryptHasher implements PasswordHasher using bcrypt.
type BcryptHasher struct{}

// CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
func (BcryptHasher) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

// GenerateFromPassword generates a bcrypt hash of the password.
func (BcryptHasher) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

// AuthService defines the interface for authentication business logic.
type AuthService interface {
	SignUp(ctx context.Context, input SignUpInput) (db.User, error)
	SignIn(ctx context.Context, input SignInInput) (AuthResult, error)
	VerifyEmail(ctx context.Context, userID int64) error
}

// SignUpInput represents the input for user registration.
type SignUpInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

// Validate checks if the sign-up input is valid.
func (i SignUpInput) Validate() error {
	// Email validation
	if i.Email == "" {
		return fmt.Errorf("%w: email is required", ErrInvalidInput)
	}
	email := strings.TrimSpace(strings.ToLower(i.Email))
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%w: invalid email format", ErrInvalidInput)
	}

	// Password validation
	if i.Password == "" {
		return fmt.Errorf("%w: password is required", ErrInvalidInput)
	}
	if len(i.Password) < MinPasswordLength {
		return fmt.Errorf("%w: password must be at least %d characters", ErrInvalidInput, MinPasswordLength)
	}

	// Name validation
	if strings.TrimSpace(i.FirstName) == "" {
		return fmt.Errorf("%w: first name is required", ErrInvalidInput)
	}
	if strings.TrimSpace(i.LastName) == "" {
		return fmt.Errorf("%w: last name is required", ErrInvalidInput)
	}

	return nil
}

// SignInInput represents the input for user login.
type SignInInput struct {
	Email      string
	Password   string
	RememberMe bool
}

// Validate checks if the sign-in input is valid.
func (i SignInInput) Validate() error {
	if i.Email == "" {
		return fmt.Errorf("%w: email is required", ErrInvalidInput)
	}
	if i.Password == "" {
		return fmt.Errorf("%w: password is required", ErrInvalidInput)
	}
	return nil
}

// AuthResult contains the result of a successful authentication.
type AuthResult struct {
	User       db.User
	RememberMe bool
}

type authService struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
	clock    Clock
	hasher   PasswordHasher
}

// NewAuthService creates a new AuthService with the given repositories.
func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository) AuthService {
	return &authService{
		authRepo: authRepo,
		userRepo: userRepo,
		clock:    RealClock{},
		hasher:   BcryptHasher{},
	}
}

func (s *authService) SignUp(ctx context.Context, input SignUpInput) (db.User, error) {
	// Validate input
	if err := input.Validate(); err != nil {
		return db.User{}, err
	}

	// Normalize email
	email := strings.TrimSpace(strings.ToLower(input.Email))

	// Check if email already exists
	_, err := s.authRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return db.User{}, ErrEmailExists
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return db.User{}, fmt.Errorf("failed to check email existence: %w", err)
	}

	// Hash password
	passwordHash, err := s.hasher.GenerateFromPassword([]byte(input.Password), BcryptCost)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.userRepo.Create(ctx, db.CreateUserParams{
		Email:     email,
		FirstName: strings.TrimSpace(input.FirstName),
		LastName:  strings.TrimSpace(input.LastName),
		Role:      db.UserRoleEntrant,
	})
	if err != nil {
		return db.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	// Create auth credentials
	_, err = s.authRepo.CreateCredentials(ctx, user.ID, string(passwordHash))
	if err != nil {
		// TODO: Consider implementing transaction rollback here
		return db.User{}, fmt.Errorf("failed to create credentials: %w", err)
	}

	// TODO: Send verification email

	return user, nil
}

func (s *authService) SignIn(ctx context.Context, input SignInInput) (AuthResult, error) {
	// Validate input
	if err := input.Validate(); err != nil {
		return AuthResult{}, err
	}

	// Normalize email
	email := strings.TrimSpace(strings.ToLower(input.Email))

	// Get user
	user, err := s.authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return AuthResult{}, ErrInvalidCredentials
		}
		return AuthResult{}, fmt.Errorf("failed to get user: %w", err)
	}

	// Get credentials
	creds, err := s.authRepo.GetCredentialsByUserID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return AuthResult{}, ErrInvalidCredentials
		}
		return AuthResult{}, fmt.Errorf("failed to get credentials: %w", err)
	}

	// Check if account is locked
	locked, err := s.authRepo.IsAccountLocked(ctx, user.ID)
	if err != nil {
		return AuthResult{}, fmt.Errorf("failed to check account lock status: %w", err)
	}
	if locked {
		return AuthResult{}, ErrAccountLocked
	}

	// Verify password
	err = s.hasher.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(input.Password))
	if err != nil {
		// Increment failed attempts
		if incrementErr := s.authRepo.IncrementFailedAttempts(ctx, user.ID); incrementErr != nil {
			// Log error but continue
		}

		// Lock account if max attempts reached
		if creds.FailedLoginAttempts+1 >= MaxLoginAttempts {
			lockUntil := s.clock.Now().Add(AccountLockoutDuration)
			if lockErr := s.authRepo.LockAccount(ctx, user.ID, lockUntil); lockErr != nil {
				// Log error but continue
			}
			return AuthResult{}, ErrAccountLocked
		}

		return AuthResult{}, ErrInvalidCredentials
	}

	// Check email verification
	if !creds.EmailVerifiedAt.Valid {
		return AuthResult{}, ErrEmailNotVerified
	}

	// Update last login and reset failed attempts
	if err := s.authRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		// Log error but don't fail the login
	}

	return AuthResult{
		User:       user,
		RememberMe: input.RememberMe,
	}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, userID int64) error {
	return s.authRepo.VerifyEmail(ctx, userID)
}
