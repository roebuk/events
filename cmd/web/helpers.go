package main

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/a-h/templ"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	//nolint:errcheck // Best effort write in error handler
	w.Write([]byte(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Server Error</title>
			<style>
				body { font-family: sans-serif; background: #f8f8f8; color: #333; padding: 2em; }
				.container { max-width: 600px; margin: auto; background: #fff; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.1); padding: 2em; }
				h1 { color: #c00; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>500 - Server Error</h1>
				<p>Sorry, something went wrong on our end.</p>
			</div>
		</body>
		</html>
	`))
}

//nolint:unparam // status parameter kept for future flexibility with different HTTP status codes
func (app *application) render(ctx context.Context, w http.ResponseWriter, status int, component templ.Component) {
	w.WriteHeader(status)

	if err := component.Render(ctx, w); err != nil {
		app.logger.Error("failed to render component", "error", err)
	}
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Flash message types
const (
	FlashSuccess = "success"
	FlashError   = "error"
	FlashInfo    = "info"
	FlashWarning = "warning"
)

// addFlash adds a flash message to the session.
func (app *application) addFlash(r *http.Request, messageType, message string) {
	app.sessionManager.Put(r.Context(), "flash_"+messageType, message)
}

// getFlash retrieves and removes a flash message from the session.
func (app *application) getFlash(r *http.Request, messageType string) string {
	return app.sessionManager.PopString(r.Context(), "flash_"+messageType)
}

// getAllFlashes retrieves all flash messages and returns them in a map.
func (app *application) getAllFlashes(r *http.Request) map[string]string {
	flashes := make(map[string]string)

	if msg := app.getFlash(r, FlashSuccess); msg != "" {
		flashes[FlashSuccess] = msg
	}
	if msg := app.getFlash(r, FlashError); msg != "" {
		flashes[FlashError] = msg
	}
	if msg := app.getFlash(r, FlashInfo); msg != "" {
		flashes[FlashInfo] = msg
	}
	if msg := app.getFlash(r, FlashWarning); msg != "" {
		flashes[FlashWarning] = msg
	}

	return flashes
}

// isAuthenticated returns true if the user is authenticated.
func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "userID")
}

// getUserID retrieves the authenticated user's ID from the session.
func (app *application) getUserID(r *http.Request) int64 {
	return app.sessionManager.GetInt64(r.Context(), "userID")
}
