package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New()

	mux.Handle("/", dynamic.ThenFunc(app.home))
	mux.Handle("GET /events/{slug}", dynamic.ThenFunc(app.eventView))

	// Authentication routes
	mux.Handle("GET /auth/sign-in", dynamic.ThenFunc(app.signInView))
	mux.Handle("POST /auth/sign-in", dynamic.ThenFunc(app.signInPost))
	mux.Handle("GET /auth/sign-up", dynamic.ThenFunc(app.signUpView))
	mux.Handle("POST /auth/sign-up", dynamic.ThenFunc(app.signUpPost))

	mux.HandleFunc("GET /insert", app.adminCreatePost)
	mux.HandleFunc("GET /insert-user", app.adminCreateUser)
	// mux.Handle("GET /admin/create", app. )
	// mux.HandleFunc("GET /privacy", app.privacyView)
	standard := alice.New(app.logRequest, commonHeaders)

	return standard.Then(mux)
}
