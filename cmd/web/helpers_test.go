package main

import (
	"firecrest-go/ui/templates/auth"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerError(t *testing.T) {
	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	app.serverError(w, req, assert.AnError)

	res := w.Result()
	defer res.Body.Close()

	// Assert status code
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	// Assert content type
	assert.Equal(t, "text/html; charset=utf-8", res.Header.Get("Content-Type"))

	// Assert body contains error page
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	bodyStr := string(body)

	assert.Contains(t, bodyStr, "500 - Server Error")
	assert.Contains(t, bodyStr, "Sorry, something went wrong on our end.")
}

func TestRender(t *testing.T) {
	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	w := httptest.NewRecorder()

	// Use a simple template component for testing
	component := auth.SignIn()

	app.render(w, http.StatusOK, component)

	res := w.Result()
	defer res.Body.Close()

	// Assert status code
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Assert body is not empty
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.NotEmpty(t, body)
}

func TestClientError(t *testing.T) {
	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	tests := []struct {
		name       string
		statusCode int
		expected   string
	}{
		{
			name:       "bad request",
			statusCode: http.StatusBadRequest,
			expected:   "Bad Request",
		},
		{
			name:       "unauthorized",
			statusCode: http.StatusUnauthorized,
			expected:   "Unauthorized",
		},
		{
			name:       "forbidden",
			statusCode: http.StatusForbidden,
			expected:   "Forbidden",
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
			expected:   "Not Found",
		},
		{
			name:       "method not allowed",
			statusCode: http.StatusMethodNotAllowed,
			expected:   "Method Not Allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			app.clientError(w, tt.statusCode)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.statusCode, res.StatusCode)

			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Contains(t, strings.TrimSpace(string(body)), tt.expected)
		})
	}
}

func TestNotFound(t *testing.T) {
	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	w := httptest.NewRecorder()

	app.notFound(w)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Contains(t, strings.TrimSpace(string(body)), "Not Found")
}
