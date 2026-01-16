package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"firecrest/db"
	"firecrest/internal/repository"
	"firecrest/internal/service"
)

// mockEventService implements service.EventService for testing.
type mockEventService struct {
	listEventsFunc  func(ctx context.Context) ([]db.Event, error)
	getEventFunc    func(ctx context.Context, slug string) (db.Event, error)
	createEventFunc func(ctx context.Context, input service.CreateEventInput) (db.Event, error)
}

func (m *mockEventService) ListEvents(ctx context.Context) ([]db.Event, error) {
	if m.listEventsFunc != nil {
		return m.listEventsFunc(ctx)
	}
	return nil, nil
}

func (m *mockEventService) GetEvent(ctx context.Context, slug string) (db.Event, error) {
	if m.getEventFunc != nil {
		return m.getEventFunc(ctx, slug)
	}
	return db.Event{}, nil
}

func (m *mockEventService) CreateEvent(ctx context.Context, input service.CreateEventInput) (db.Event, error) {
	if m.createEventFunc != nil {
		return m.createEventFunc(ctx, input)
	}
	return db.Event{}, nil
}

// mockUserService implements service.UserService for testing.
type mockUserService struct {
	getUserFunc    func(ctx context.Context, id int64) (db.User, error)
	createUserFunc func(ctx context.Context, input service.CreateUserInput) (db.User, error)
}

func (m *mockUserService) GetUser(ctx context.Context, id int64) (db.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(ctx, id)
	}
	return db.User{}, nil
}

func (m *mockUserService) CreateUser(ctx context.Context, input service.CreateUserInput) (db.User, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, input)
	}
	return db.User{}, nil
}

func newTestApplication(eventSvc service.EventService, userSvc service.UserService) *application {
	return &application{
		logger:       slog.New(slog.NewTextHandler(io.Discard, nil)),
		eventService: eventSvc,
		userService:  userSvc,
	}
}

func TestHome(t *testing.T) {
	t.Run("returns 200 and renders events", func(t *testing.T) {
		mockEventSvc := &mockEventService{
			listEventsFunc: func(ctx context.Context) ([]db.Event, error) {
				return []db.Event{
					{ID: 1, Name: "Test Event 1", Slug: "test-event-1"},
					{ID: 2, Name: "Test Event 2", Slug: "test-event-2"},
				}, nil
			},
		}

		app := newTestApplication(mockEventSvc, &mockUserService{})

		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rr := httptest.NewRecorder()

		app.home(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("returns 200 with empty events list", func(t *testing.T) {
		mockEventSvc := &mockEventService{
			listEventsFunc: func(ctx context.Context) ([]db.Event, error) {
				return []db.Event{}, nil
			},
		}

		app := newTestApplication(mockEventSvc, &mockUserService{})

		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rr := httptest.NewRecorder()

		app.home(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("returns 500 on service error", func(t *testing.T) {
		mockEventSvc := &mockEventService{
			listEventsFunc: func(ctx context.Context) ([]db.Event, error) {
				return nil, errors.New("database connection failed")
			},
		}

		app := newTestApplication(mockEventSvc, &mockUserService{})

		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rr := httptest.NewRecorder()

		app.home(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}

func TestEventView(t *testing.T) {
	t.Run("returns 200 for valid event", func(t *testing.T) {
		mockEventSvc := &mockEventService{
			getEventFunc: func(ctx context.Context, slug string) (db.Event, error) {
				if slug == "test-event" {
					return db.Event{
						ID:   1,
						Name: "Test Event",
						Slug: "test-event",
					}, nil
				}
				return db.Event{}, repository.ErrNotFound
			},
		}

		app := newTestApplication(mockEventSvc, &mockUserService{})

		req := httptest.NewRequest(http.MethodGet, "/events/test-event", http.NoBody)
		req.SetPathValue("slug", "test-event")
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("returns 404 for non-existent event", func(t *testing.T) {
		mockEventSvc := &mockEventService{
			getEventFunc: func(ctx context.Context, slug string) (db.Event, error) {
				return db.Event{}, repository.ErrNotFound
			},
		}

		app := newTestApplication(mockEventSvc, &mockUserService{})

		req := httptest.NewRequest(http.MethodGet, "/events/non-existent", http.NoBody)
		req.SetPathValue("slug", "non-existent")
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("returns 400 for empty slug", func(t *testing.T) {
		mockEventSvc := &mockEventService{
			getEventFunc: func(ctx context.Context, slug string) (db.Event, error) {
				return db.Event{}, service.ErrInvalidInput
			},
		}
		app := newTestApplication(mockEventSvc, &mockUserService{})

		req := httptest.NewRequest(http.MethodGet, "/events/", http.NoBody)
		req.SetPathValue("slug", "")
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("returns 400 for slug exceeding 100 characters", func(t *testing.T) {
		mockEventSvc := &mockEventService{
			getEventFunc: func(ctx context.Context, slug string) (db.Event, error) {
				return db.Event{}, service.ErrInvalidInput
			},
		}
		app := newTestApplication(mockEventSvc, &mockUserService{})
		longSlug := strings.Repeat("a", 101)

		req := httptest.NewRequest(http.MethodGet, "/events/"+longSlug, http.NoBody)
		req.SetPathValue("slug", longSlug)
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("returns 500 on service error", func(t *testing.T) {
		mockEventSvc := &mockEventService{
			getEventFunc: func(ctx context.Context, slug string) (db.Event, error) {
				return db.Event{}, errors.New("database connection failed")
			},
		}

		app := newTestApplication(mockEventSvc, &mockUserService{})

		req := httptest.NewRequest(http.MethodGet, "/events/test-event", http.NoBody)
		req.SetPathValue("slug", "test-event")
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}
