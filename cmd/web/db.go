package main

import (
	"context"
	"firecrest-go/db"
)

// DB defines the interface for database operations
type DB interface {
	ListEvents(ctx context.Context) ([]db.Event, error)
	GetEvent(ctx context.Context, slug string) (db.Event, error)
	CreateEvent(ctx context.Context, params db.CreateEventParams) (db.Event, error)
	CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error)
}
