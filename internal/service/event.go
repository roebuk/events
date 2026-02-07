package service

import (
	"context"
	"errors"
	"fmt"

	"firecrest/db"
	"firecrest/internal/repository"
)

// ErrInvalidInput is returned when input validation fails.
var ErrInvalidInput = errors.New("invalid input")

// EventService defines the interface for event business logic.
type EventService interface {
	ListEvents(ctx context.Context) ([]db.Event, error)
	GetEvent(ctx context.Context, slug string) (db.Event, error)
	CreateEvent(ctx context.Context, input CreateEventInput) (db.Event, error)
}

// CreateEventInput represents the input for creating an event.
type CreateEventInput struct {
	OrganisationID int64
	Name           string
	Slug           string
	Year           int32
}

// Validate checks if the input is valid.
func (i CreateEventInput) Validate() error {
	if i.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if i.Slug == "" {
		return fmt.Errorf("%w: slug is required", ErrInvalidInput)
	}
	if len(i.Slug) > 100 {
		return fmt.Errorf("%w: slug must be 100 characters or less", ErrInvalidInput)
	}
	if i.OrganisationID <= 0 {
		return fmt.Errorf("%w: organisation_id must be positive", ErrInvalidInput)
	}
	if i.Year < 2025 {
		return fmt.Errorf("%w: year must be 2025 or later", ErrInvalidInput)
	}
	return nil
}

type eventService struct {
	eventRepo repository.EventRepository
}

// NewEventService creates a new EventService with the given repository.
func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) ListEvents(ctx context.Context) ([]db.Event, error) {
	return s.eventRepo.List(ctx)
}

func (s *eventService) GetEvent(ctx context.Context, slug string) (db.Event, error) {
	if slug == "" || len(slug) > 100 {
		return db.Event{}, fmt.Errorf("%w: invalid slug", ErrInvalidInput)
	}
	return s.eventRepo.GetBySlug(ctx, slug)
}

func (s *eventService) CreateEvent(ctx context.Context, input CreateEventInput) (db.Event, error) {
	if err := input.Validate(); err != nil {
		return db.Event{}, err
	}

	return s.eventRepo.Create(ctx, db.CreateEventParams{
		OrganisationID: input.OrganisationID,
		Name:           input.Name,
		Slug:           input.Slug,
		Year:           input.Year,
	})
}
