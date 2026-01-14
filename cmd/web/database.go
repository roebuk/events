package main

import (
	"context"

	"firecrest/db"
)

// Database defines the interface for database operations used by handlers.
// This allows for easy mocking in tests.
type Database interface {
	ListEvents(ctx context.Context) ([]db.Event, error)
	GetEvent(ctx context.Context, slug string) (db.Event, error)
	CreateEvent(ctx context.Context, arg db.CreateEventParams) (db.Event, error)
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
}
