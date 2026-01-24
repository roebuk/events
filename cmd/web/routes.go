package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /health", app.health)

	// Protects against CSRF by checking Sec-Fetch-Site header
	// https://www.alexedwards.net/blog/preventing-csrf-in-go
	cop := http.NewCrossOriginProtection()
	cop.AddTrustedOrigin("http://localhost:8080")

	// Middleware chains
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.loadUser)
	authRequired := dynamic.Append(app.requireAuth)
	guestOnly := dynamic.Append(app.redirectIfAuth)

	// Public routes
	mux.Handle("GET /", dynamic.ThenFunc(app.home))
	mux.Handle("GET /events/{slug}", dynamic.ThenFunc(app.eventView))

	// Authentication routes (guest only)
	mux.Handle("GET /auth/sign-in", guestOnly.ThenFunc(app.signInView))
	mux.Handle("POST /auth/sign-in", guestOnly.ThenFunc(app.signInPost))
	mux.Handle("GET /auth/sign-up", guestOnly.ThenFunc(app.signUpView))
	mux.Handle("POST /auth/sign-up", guestOnly.ThenFunc(app.signUpPost))
	mux.Handle("GET /auth/verify-email", dynamic.ThenFunc(app.verifyEmail))

	// Sign out (authenticated only)
	mux.Handle("POST /auth/sign-out", authRequired.ThenFunc(app.signOut))

	// Admin routes (temporary - should be removed in production)
	mux.HandleFunc("GET /insert", app.adminCreatePost)
	mux.HandleFunc("GET /insert-user", app.adminCreateUser)

	// Apply standard middleware + Cross-Origin Protection
	standard := alice.New(app.logRequest, commonHeaders)

	return standard.Then(cop.Handler(mux))
}
