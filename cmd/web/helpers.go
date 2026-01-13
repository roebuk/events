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

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
