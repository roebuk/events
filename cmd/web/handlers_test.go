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

	"github.com/jackc/pgx/v5"

	"firecrest/db"
)

// mockDB implements the Database interface for testing.
type mockDB struct {
	listEventsFunc  func(ctx context.Context) ([]db.Event, error)
	getEventFunc    func(ctx context.Context, slug string) (db.Event, error)
	createEventFunc func(ctx context.Context, arg db.CreateEventParams) (db.Event, error)
	createUserFunc  func(ctx context.Context, arg db.CreateUserParams) (db.User, error)
}

func (m *mockDB) ListEvents(ctx context.Context) ([]db.Event, error) {
	if m.listEventsFunc != nil {
		return m.listEventsFunc(ctx)
	}
	return nil, nil
}

func (m *mockDB) GetEvent(ctx context.Context, slug string) (db.Event, error) {
	if m.getEventFunc != nil {
		return m.getEventFunc(ctx, slug)
	}
	return db.Event{}, nil
}

func (m *mockDB) CreateEvent(ctx context.Context, arg db.CreateEventParams) (db.Event, error) {
	if m.createEventFunc != nil {
		return m.createEventFunc(ctx, arg)
	}
	return db.Event{}, nil
}

func (m *mockDB) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	if m.createUserFunc != nil {
		return m.createUserFunc(ctx, arg)
	}
	return db.User{}, nil
}

func newTestApplication(database Database) *application {
	return &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		db:     database,
	}
}

func TestHome(t *testing.T) {
	t.Run("returns 200 and renders events", func(t *testing.T) {
		mock := &mockDB{
			listEventsFunc: func(ctx context.Context) ([]db.Event, error) {
				return []db.Event{
					{ID: 1, Name: "Test Event 1", Slug: "test-event-1"},
					{ID: 2, Name: "Test Event 2", Slug: "test-event-2"},
				}, nil
			},
		}

		app := newTestApplication(mock)

		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rr := httptest.NewRecorder()

		app.home(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("returns 200 with empty events list", func(t *testing.T) {
		mock := &mockDB{
			listEventsFunc: func(ctx context.Context) ([]db.Event, error) {
				return []db.Event{}, nil
			},
		}

		app := newTestApplication(mock)

		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rr := httptest.NewRecorder()

		app.home(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("returns 500 on database error", func(t *testing.T) {
		mock := &mockDB{
			listEventsFunc: func(ctx context.Context) ([]db.Event, error) {
				return nil, errors.New("database connection failed")
			},
		}

		app := newTestApplication(mock)

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
		mock := &mockDB{
			getEventFunc: func(ctx context.Context, slug string) (db.Event, error) {
				if slug == "test-event" {
					return db.Event{
						ID:   1,
						Name: "Test Event",
						Slug: "test-event",
					}, nil
				}
				return db.Event{}, pgx.ErrNoRows
			},
		}

		app := newTestApplication(mock)

		req := httptest.NewRequest(http.MethodGet, "/events/test-event", http.NoBody)
		req.SetPathValue("slug", "test-event")
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("returns 404 for non-existent event", func(t *testing.T) {
		mock := &mockDB{
			getEventFunc: func(ctx context.Context, slug string) (db.Event, error) {
				return db.Event{}, pgx.ErrNoRows
			},
		}

		app := newTestApplication(mock)

		req := httptest.NewRequest(http.MethodGet, "/events/non-existent", http.NoBody)
		req.SetPathValue("slug", "non-existent")
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, rr.Code)
		}
	})

	t.Run("returns 400 for empty slug", func(t *testing.T) {
		mock := &mockDB{}
		app := newTestApplication(mock)

		req := httptest.NewRequest(http.MethodGet, "/events/", http.NoBody)
		req.SetPathValue("slug", "")
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("returns 400 for slug exceeding 100 characters", func(t *testing.T) {
		mock := &mockDB{}
		app := newTestApplication(mock)
		longSlug := strings.Repeat("a", 101)

		req := httptest.NewRequest(http.MethodGet, "/events/"+longSlug, http.NoBody)
		req.SetPathValue("slug", longSlug)
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("returns 500 on database error", func(t *testing.T) {
		mock := &mockDB{
			getEventFunc: func(ctx context.Context, slug string) (db.Event, error) {
				return db.Event{}, errors.New("database connection failed")
			},
		}

		app := newTestApplication(mock)

		req := httptest.NewRequest(http.MethodGet, "/events/test-event", http.NoBody)
		req.SetPathValue("slug", "test-event")
		rr := httptest.NewRecorder()

		app.eventView(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}
