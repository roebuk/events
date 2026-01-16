package service

import (
	"context"
	"errors"
	"testing"

	"firecrest/db"
	"firecrest/internal/repository"
)

// mockEventRepository implements repository.EventRepository for testing.
type mockEventRepository struct {
	listFunc      func(ctx context.Context) ([]db.Event, error)
	getBySlugFunc func(ctx context.Context, slug string) (db.Event, error)
	createFunc    func(ctx context.Context, params db.CreateEventParams) (db.Event, error)
}

func (m *mockEventRepository) List(ctx context.Context) ([]db.Event, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return nil, nil
}

func (m *mockEventRepository) GetBySlug(ctx context.Context, slug string) (db.Event, error) {
	if m.getBySlugFunc != nil {
		return m.getBySlugFunc(ctx, slug)
	}
	return db.Event{}, nil
}

func (m *mockEventRepository) Create(ctx context.Context, params db.CreateEventParams) (db.Event, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, params)
	}
	return db.Event{}, nil
}

func TestEventService_ListEvents(t *testing.T) {
	t.Run("returns events from repository", func(t *testing.T) {
		expected := []db.Event{
			{ID: 1, Name: "Event 1", Slug: "event-1"},
			{ID: 2, Name: "Event 2", Slug: "event-2"},
		}

		repo := &mockEventRepository{
			listFunc: func(ctx context.Context) ([]db.Event, error) {
				return expected, nil
			},
		}

		svc := NewEventService(repo)
		events, err := svc.ListEvents(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(events) != len(expected) {
			t.Errorf("expected %d events, got %d", len(expected), len(events))
		}
	})

	t.Run("propagates repository errors", func(t *testing.T) {
		repo := &mockEventRepository{
			listFunc: func(ctx context.Context) ([]db.Event, error) {
				return nil, errors.New("database error")
			},
		}

		svc := NewEventService(repo)
		_, err := svc.ListEvents(context.Background())

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestEventService_GetEvent(t *testing.T) {
	t.Run("returns event for valid slug", func(t *testing.T) {
		expected := db.Event{ID: 1, Name: "Test Event", Slug: "test-event"}

		repo := &mockEventRepository{
			getBySlugFunc: func(ctx context.Context, slug string) (db.Event, error) {
				if slug == "test-event" {
					return expected, nil
				}
				return db.Event{}, repository.ErrNotFound
			},
		}

		svc := NewEventService(repo)
		event, err := svc.GetEvent(context.Background(), "test-event")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if event.ID != expected.ID {
			t.Errorf("expected event ID %d, got %d", expected.ID, event.ID)
		}
	})

	t.Run("returns ErrInvalidInput for empty slug", func(t *testing.T) {
		repo := &mockEventRepository{}
		svc := NewEventService(repo)

		_, err := svc.GetEvent(context.Background(), "")

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrInvalidInput for slug exceeding 100 characters", func(t *testing.T) {
		repo := &mockEventRepository{}
		svc := NewEventService(repo)

		longSlug := make([]byte, 101)
		for i := range longSlug {
			longSlug[i] = 'a'
		}

		_, err := svc.GetEvent(context.Background(), string(longSlug))

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrNotFound for non-existent event", func(t *testing.T) {
		repo := &mockEventRepository{
			getBySlugFunc: func(ctx context.Context, slug string) (db.Event, error) {
				return db.Event{}, repository.ErrNotFound
			},
		}

		svc := NewEventService(repo)
		_, err := svc.GetEvent(context.Background(), "non-existent")

		if !errors.Is(err, repository.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestEventService_CreateEvent(t *testing.T) {
	t.Run("creates event with valid input", func(t *testing.T) {
		expected := db.Event{ID: 1, Name: "New Event", Slug: "new-event", OrganisationID: 1}

		repo := &mockEventRepository{
			createFunc: func(ctx context.Context, params db.CreateEventParams) (db.Event, error) {
				return expected, nil
			},
		}

		svc := NewEventService(repo)
		event, err := svc.CreateEvent(context.Background(), CreateEventInput{
			OrganisationID: 1,
			Name:           "New Event",
			Slug:           "new-event",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if event.ID != expected.ID {
			t.Errorf("expected event ID %d, got %d", expected.ID, event.ID)
		}
	})

	t.Run("returns ErrInvalidInput for missing name", func(t *testing.T) {
		repo := &mockEventRepository{}
		svc := NewEventService(repo)

		_, err := svc.CreateEvent(context.Background(), CreateEventInput{
			OrganisationID: 1,
			Name:           "",
			Slug:           "new-event",
		})

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrInvalidInput for missing slug", func(t *testing.T) {
		repo := &mockEventRepository{}
		svc := NewEventService(repo)

		_, err := svc.CreateEvent(context.Background(), CreateEventInput{
			OrganisationID: 1,
			Name:           "New Event",
			Slug:           "",
		})

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("returns ErrInvalidInput for invalid organisation_id", func(t *testing.T) {
		repo := &mockEventRepository{}
		svc := NewEventService(repo)

		_, err := svc.CreateEvent(context.Background(), CreateEventInput{
			OrganisationID: 0,
			Name:           "New Event",
			Slug:           "new-event",
		})

		if !errors.Is(err, ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("propagates repository errors", func(t *testing.T) {
		repo := &mockEventRepository{
			createFunc: func(ctx context.Context, params db.CreateEventParams) (db.Event, error) {
				return db.Event{}, errors.New("database error")
			},
		}

		svc := NewEventService(repo)
		_, err := svc.CreateEvent(context.Background(), CreateEventInput{
			OrganisationID: 1,
			Name:           "New Event",
			Slug:           "new-event",
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}
