package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("GET /events/{slug}", app.eventView)

	// Authentication routes
	mux.HandleFunc("GET /auth/sign-in", app.signInView)
	mux.HandleFunc("POST /auth/sign-in", app.signInPost)
	mux.HandleFunc("GET /auth/sign-up", app.signUpView)
	mux.HandleFunc("POST /auth/sign-up", app.signUpPost)

	// mux.HandleFunc("GET /privacy", app.privacyView)

	return mux
}
