package main

import (
	"context"
	"errors"
	"firecrest-go/db"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// mockDB is a mock implementation of dbInterface
type mockDB struct {
	listEventsFunc func(ctx context.Context) ([]db.Event, error)
}

func (m *mockDB) ListEvents(ctx context.Context) ([]db.Event, error) {
	if m.listEventsFunc != nil {
		return m.listEventsFunc(ctx)
	}
	return []db.Event{}, nil
}

func (m *mockDB) GetEvent(ctx context.Context, slug string) (db.Event, error) {
	return db.Event{}, nil
}

func (m *mockDB) CreateEvent(ctx context.Context, arg db.CreateEventParams) (db.Event, error) {
	return db.Event{}, nil
}

func (m *mockDB) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return db.User{}, nil
}

func TestHome(t *testing.T) {
	tests := []struct {
		name           string
		listEventsFunc func(ctx context.Context) ([]db.Event, error)
		wantStatus     int
	}{
		{
			name: "successfully lists events",
			listEventsFunc: func(ctx context.Context) ([]db.Event, error) {
				return []db.Event{
					{
						ID:   1,
						Name: "Test Event",
						Slug: "test-event",
					},
				}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles empty events list",
			listEventsFunc: func(ctx context.Context) ([]db.Event, error) {
				return []db.Event{}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "handles database error",
			listEventsFunc: func(ctx context.Context) ([]db.Event, error) {
				return nil, errors.New("database error")
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test application with mocked database
			mock := &mockDB{
				listEventsFunc: tt.listEventsFunc,
			}

			app := &application{
				logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
				db:     mock,
			}

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			tr := httptest.NewRecorder()

			app.home(tr, req)

			// Check the status code
			if tr.Code != tt.wantStatus {
				t.Errorf("home() status = %v, want %v", tr.Code, tt.wantStatus)
			}

		})
	}
}
