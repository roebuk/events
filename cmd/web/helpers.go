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

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, component templ.Component) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "text/html")

	component.Render(context.Background(), w)
}
