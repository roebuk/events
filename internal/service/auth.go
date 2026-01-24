package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"firecrest/db"
	"firecrest/internal/repository"
)

// Authentication errors
var (
	ErrInvalidCredentials      = errors.New("invalid email or password")
	ErrEmailNotVerified        = errors.New("email address not verified")
	ErrAccountLocked           = errors.New("account is locked due to too many failed login attempts")
	ErrEmailExists             = errors.New("email address already registered")
	ErrInvalidVerificationCode = errors.New("invalid or expired verification code")
)

// Authentication constants
const (
	BcryptCost                  = 12
	MaxLoginAttempts            = 5
	AccountLockoutDuration      = 15 * time.Minute
	MinPasswordLength           = 8
	VerificationTokenExpiry     = 24 * time.Hour
	verificationTokenSecret     = "email-verification-secret" // TODO: Move to environment variable
	verificationTokenSeparator  = "."
)

// AuthService defines the interface for authentication business logic.
type AuthService interface {
	SignUp(ctx context.Context, input SignUpInput) (SignUpResult, error)
	SignIn(ctx context.Context, input SignInInput) (AuthResult, error)
	VerifyEmail(ctx context.Context, userID int64) error
	VerifyEmailByToken(ctx context.Context, token string) error
}

// SignUpResult contains the result of a successful sign-up.
type SignUpResult struct {
	User             db.User
	VerificationCode string
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
}

// NewAuthService creates a new AuthService with the given repositories.
func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository) AuthService {
	return &authService{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s *authService) SignUp(ctx context.Context, input SignUpInput) (SignUpResult, error) {
	// Validate input
	if err := input.Validate(); err != nil {
		return SignUpResult{}, err
	}

	// Normalize email
	email := strings.TrimSpace(strings.ToLower(input.Email))

	// Check if email already exists
	_, err := s.authRepo.GetUserByEmail(ctx, email)
	if err == nil {
		return SignUpResult{}, ErrEmailExists
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return SignUpResult{}, fmt.Errorf("failed to check email existence: %w", err)
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), BcryptCost)
	if err != nil {
		return SignUpResult{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user, err := s.userRepo.Create(ctx, db.CreateUserParams{
		Email:     email,
		FirstName: strings.TrimSpace(input.FirstName),
		LastName:  strings.TrimSpace(input.LastName),
		Role:      db.UserRoleEntrant, // Default role
	})
	if err != nil {
		return SignUpResult{}, fmt.Errorf("failed to create user: %w", err)
	}

	// Create auth credentials
	_, err = s.authRepo.CreateCredentials(ctx, user.ID, string(passwordHash))
	if err != nil {
		// TODO: Consider implementing transaction rollback here
		return SignUpResult{}, fmt.Errorf("failed to create credentials: %w", err)
	}

	// Generate verification code
	verificationCode := generateVerificationToken(user.ID)

	// TODO: Send verification email in production

	return SignUpResult{
		User:             user,
		VerificationCode: verificationCode,
	}, nil
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
	err = bcrypt.CompareHashAndPassword([]byte(creds.PasswordHash), []byte(input.Password))
	if err != nil {
		// Increment failed attempts
		if incrementErr := s.authRepo.IncrementFailedAttempts(ctx, user.ID); incrementErr != nil {
			// Log error but continue
		}

		// Lock account if max attempts reached
		if creds.FailedLoginAttempts+1 >= MaxLoginAttempts {
			lockUntil := time.Now().Add(AccountLockoutDuration)
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

func (s *authService) VerifyEmailByToken(ctx context.Context, token string) error {
	userID, err := validateVerificationToken(token)
	if err != nil {
		return ErrInvalidVerificationCode
	}

	return s.authRepo.VerifyEmail(ctx, userID)
}

// generateVerificationToken creates a signed token containing the user ID and expiry time.
// Format: base64(userID.expiryTimestamp).signature
func generateVerificationToken(userID int64) string {
	expiry := time.Now().Add(VerificationTokenExpiry).Unix()
	payload := fmt.Sprintf("%d%s%d", userID, verificationTokenSeparator, expiry)

	// Create HMAC signature
	h := hmac.New(sha256.New, []byte(verificationTokenSecret))
	h.Write([]byte(payload))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Encode payload
	encodedPayload := base64.URLEncoding.EncodeToString([]byte(payload))

	return encodedPayload + verificationTokenSeparator + signature
}

// validateVerificationToken validates the token and returns the user ID if valid.
func validateVerificationToken(token string) (int64, error) {
	parts := strings.Split(token, verificationTokenSeparator)
	if len(parts) != 2 {
		return 0, errors.New("invalid token format")
	}

	encodedPayload, providedSignature := parts[0], parts[1]

	// Decode payload
	payloadBytes, err := base64.URLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return 0, errors.New("invalid token encoding")
	}
	payload := string(payloadBytes)

	// Verify signature
	h := hmac.New(sha256.New, []byte(verificationTokenSecret))
	h.Write(payloadBytes)
	expectedSignature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(providedSignature), []byte(expectedSignature)) {
		return 0, errors.New("invalid token signature")
	}

	// Parse payload
	payloadParts := strings.Split(payload, verificationTokenSeparator)
	if len(payloadParts) != 2 {
		return 0, errors.New("invalid payload format")
	}

	userID, err := strconv.ParseInt(payloadParts[0], 10, 64)
	if err != nil {
		return 0, errors.New("invalid user ID in token")
	}

	expiry, err := strconv.ParseInt(payloadParts[1], 10, 64)
	if err != nil {
		return 0, errors.New("invalid expiry in token")
	}

	// Check expiry
	if time.Now().Unix() > expiry {
		return 0, errors.New("token expired")
	}

	return userID, nil
}
