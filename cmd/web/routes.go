package main

import (
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /health", app.health)

	// CSRF Protection Middleware
	csrfKey := []byte(os.Getenv("CSRF_KEY"))
	csrfMiddleware := csrf.Protect(
		csrfKey,
		csrf.Secure(false), // Set to true in production with HTTPS
		csrf.Path("/"),
		csrf.SameSite(csrf.SameSiteLaxMode),
	)

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

	// Sign out (authenticated only)
	mux.Handle("POST /auth/sign-out", authRequired.ThenFunc(app.signOut))

	// Admin routes (temporary - should be removed in production)
	mux.HandleFunc("GET /insert", app.adminCreatePost)
	mux.HandleFunc("GET /insert-user", app.adminCreateUser)

	// Apply standard middleware + CSRF
	standard := alice.New(app.logRequest, commonHeaders)

	return standard.Then(csrfMiddleware(mux))
}
