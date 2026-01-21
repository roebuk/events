package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"firecrest/db"
)

// AuthRepository defines the interface for authentication data access.
type AuthRepository interface {
	// User operations
	GetUserByEmail(ctx context.Context, email string) (db.User, error)

	// Credentials operations
	CreateCredentials(ctx context.Context, userID int64, passwordHash string) (db.AuthCredential, error)
	GetCredentialsByEmail(ctx context.Context, email string) (db.AuthCredential, error)
	GetCredentialsByUserID(ctx context.Context, userID int64) (db.AuthCredential, error)

	// Login tracking
	UpdateLastLogin(ctx context.Context, userID int64) error
	IncrementFailedAttempts(ctx context.Context, userID int64) error

	// Account locking
	LockAccount(ctx context.Context, userID int64, lockUntil time.Time) error
	IsAccountLocked(ctx context.Context, userID int64) (bool, error)

	// Email verification
	VerifyEmail(ctx context.Context, userID int64) error
}

type authRepository struct {
	queries *db.Queries
}

// NewAuthRepository creates a new AuthRepository backed by the given queries.
func NewAuthRepository(queries *db.Queries) AuthRepository {
	return &authRepository{queries: queries}
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (r *authRepository) CreateCredentials(ctx context.Context, userID int64, passwordHash string) (db.AuthCredential, error) {
	return r.queries.CreateAuthCredentials(ctx, db.CreateAuthCredentialsParams{
		UserID:       userID,
		PasswordHash: passwordHash,
	})
}

func (r *authRepository) GetCredentialsByEmail(ctx context.Context, email string) (db.AuthCredential, error) {
	creds, err := r.queries.GetAuthCredentialsByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.AuthCredential{}, ErrNotFound
		}
		return db.AuthCredential{}, err
	}
	return creds, nil
}

func (r *authRepository) GetCredentialsByUserID(ctx context.Context, userID int64) (db.AuthCredential, error) {
	creds, err := r.queries.GetAuthCredentialsByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.AuthCredential{}, ErrNotFound
		}
		return db.AuthCredential{}, err
	}
	return creds, nil
}

func (r *authRepository) UpdateLastLogin(ctx context.Context, userID int64) error {
	return r.queries.UpdateLastLogin(ctx, userID)
}

func (r *authRepository) IncrementFailedAttempts(ctx context.Context, userID int64) error {
	return r.queries.IncrementFailedLoginAttempts(ctx, userID)
}

func (r *authRepository) LockAccount(ctx context.Context, userID int64, lockUntil time.Time) error {
	return r.queries.LockAccount(ctx, db.LockAccountParams{
		UserID:      userID,
		LockedUntil: pgtype.Timestamptz{Time: lockUntil, Valid: true},
	})
}

func (r *authRepository) IsAccountLocked(ctx context.Context, userID int64) (bool, error) {
	return r.queries.IsAccountLocked(ctx, userID)
}

func (r *authRepository) VerifyEmail(ctx context.Context, userID int64) error {
	return r.queries.VerifyEmail(ctx, userID)
}
