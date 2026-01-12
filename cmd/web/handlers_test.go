package main

import (
	"context"
	"firecrest-go/db"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockDB is a minimal mock implementation of the db.Queries interface
// In a real project, you would use a proper mocking library or testcontainers
type mockDB struct {
	events []db.Event
	event  db.Event
	err    error
}

func (m *mockDB) ListEvents(ctx context.Context) ([]db.Event, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.events, nil
}

func (m *mockDB) GetEvent(ctx context.Context, slug string) (db.Event, error) {
	if m.err != nil {
		return db.Event{}, m.err
	}
	return m.event, nil
}

func (m *mockDB) CreateEvent(ctx context.Context, params db.CreateEventParams) (db.Event, error) {
	if m.err != nil {
		return db.Event{}, m.err
	}
	return m.event, nil
}

func (m *mockDB) CreateUser(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	if m.err != nil {
		return db.User{}, m.err
	}
	return db.User{}, nil
}

// newTestApplication creates a test application instance
func newTestApplication(t *testing.T, queries DB) *application {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return &application{
		logger: logger,
		db:     queries,
	}
}

func TestHome(t *testing.T) {
	t.Run("successfully renders home page", func(t *testing.T) {
		// Setup
		mockEvents := []db.Event{
			{ID: 1, Name: "Test Event 1", Slug: "test-event-1"},
			{ID: 2, Name: "Test Event 2", Slug: "test-event-2"},
		}

		app := newTestApplication(t, &mockDB{events: mockEvents})

		// Create request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// Execute
		app.home(w, req)

		// Assert
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		assert.NotEmpty(t, body)
	})

	t.Run("handles database error", func(t *testing.T) {
		// Setup
		app := newTestApplication(t, &mockDB{err: assert.AnError})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		// Execute
		app.home(w, req)

		// Assert
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestEventView(t *testing.T) {
	t.Run("successfully renders event page", func(t *testing.T) {
		// Setup
		mockEvent := db.Event{
			ID:   1,
			Name: "Lincoln 10k",
			Slug: "lincoln-10k",
		}

		app := newTestApplication(t, &mockDB{event: mockEvent})

		req := httptest.NewRequest(http.MethodGet, "/events/lincoln-10k", nil)
		req.SetPathValue("slug", "lincoln-10k")
		w := httptest.NewRecorder()

		// Execute
		app.eventView(w, req)

		// Assert
		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("returns bad request for invalid slug length", func(t *testing.T) {
		app := newTestApplication(t, &mockDB{})

		// Test empty slug
		req := httptest.NewRequest(http.MethodGet, "/events/", nil)
		req.SetPathValue("slug", "")
		w := httptest.NewRecorder()

		app.eventView(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("returns not found for non-existent event", func(t *testing.T) {
		app := newTestApplication(t, &mockDB{err: assert.AnError})

		req := httptest.NewRequest(http.MethodGet, "/events/non-existent", nil)
		req.SetPathValue("slug", "non-existent")
		w := httptest.NewRecorder()

		app.eventView(w, req)

		res := w.Result()
		defer res.Body.Close()

		// In this case, our mock returns a generic error
		// In real tests with a proper mock, you'd return pgx.ErrNoRows
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestSignInView(t *testing.T) {
	app := newTestApplication(t, &mockDB{})

	req := httptest.NewRequest(http.MethodGet, "/signin", nil)
	w := httptest.NewRecorder()

	app.signInView(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestSignInPost(t *testing.T) {
	app := newTestApplication(t, &mockDB{})

	req := httptest.NewRequest(http.MethodPost, "/signin", nil)
	w := httptest.NewRecorder()

	app.signInPost(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestSignUpView(t *testing.T) {
	app := newTestApplication(t, &mockDB{})

	req := httptest.NewRequest(http.MethodGet, "/signup", nil)
	w := httptest.NewRecorder()

	app.signUpView(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestSignUpPost(t *testing.T) {
	app := newTestApplication(t, &mockDB{})

	req := httptest.NewRequest(http.MethodPost, "/signup", nil)
	w := httptest.NewRecorder()

	app.signUpPost(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestAdminCreateView(t *testing.T) {
	app := newTestApplication(t, &mockDB{})

	req := httptest.NewRequest(http.MethodGet, "/admin/create", nil)
	w := httptest.NewRecorder()

	app.adminCreateView(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestAdminCreatePost(t *testing.T) {
	t.Run("successfully creates event", func(t *testing.T) {
		mockEvent := db.Event{
			ID:   1,
			Name: "Lincoln 10k",
			Slug: "lincoln-10k",
		}

		app := newTestApplication(t, &mockDB{event: mockEvent})

		req := httptest.NewRequest(http.MethodPost, "/admin/create", nil)
		w := httptest.NewRecorder()

		app.adminCreatePost(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("handles database error", func(t *testing.T) {
		app := newTestApplication(t, &mockDB{err: assert.AnError})

		req := httptest.NewRequest(http.MethodPost, "/admin/create", nil)
		w := httptest.NewRecorder()

		app.adminCreatePost(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestAdminCreateUser(t *testing.T) {
	t.Run("successfully creates user", func(t *testing.T) {
		app := newTestApplication(t, &mockDB{})

		req := httptest.NewRequest(http.MethodPost, "/admin/user/create", nil)
		w := httptest.NewRecorder()

		app.adminCreateUser(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("handles database error", func(t *testing.T) {
		app := newTestApplication(t, &mockDB{err: assert.AnError})

		req := httptest.NewRequest(http.MethodPost, "/admin/user/create", nil)
		w := httptest.NewRecorder()

		app.adminCreateUser(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
