package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"firecrest/db"
)

// EventRepository defines the interface for event data access.
type EventRepository interface {
	List(ctx context.Context) ([]db.Event, error)
	GetBySlug(ctx context.Context, slug string) (db.Event, error)
	Create(ctx context.Context, params db.CreateEventParams) (db.Event, error)
}

type eventRepository struct {
	queries *db.Queries
}

// NewEventRepository creates a new EventRepository backed by the given queries.
func NewEventRepository(queries *db.Queries) EventRepository {
	return &eventRepository{queries: queries}
}

func (r *eventRepository) List(ctx context.Context) ([]db.Event, error) {
	return r.queries.ListEvents(ctx)
}

func (r *eventRepository) GetBySlug(ctx context.Context, slug string) (db.Event, error) {
	event, err := r.queries.GetEvent(ctx, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Event{}, ErrNotFound
		}
		return db.Event{}, err
	}
	return event, nil
}

func (r *eventRepository) Create(ctx context.Context, params db.CreateEventParams) (db.Event, error) {
	return r.queries.CreateEvent(ctx, params)
}
