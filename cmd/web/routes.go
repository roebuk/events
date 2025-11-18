package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("GET /events/{slug}", app.eventView)

	mux.HandleFunc("GET /privacy", app.eventView)

	return mux
}
