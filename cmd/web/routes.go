package main

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /health", app.health)

	// CSRF Protection Middleware - using environment-specific config
	csrfMiddleware := csrf.Protect(
		[]byte(app.config.CSRF.Key),
		csrf.FieldName("csrf"),
		csrf.CookieName("csrf"),
		csrf.Secure(app.config.CSRF.SecureCookie),
		csrf.Path("/"),
		csrf.SameSite(csrf.SameSiteLaxMode),
		csrf.TrustedOrigins(app.config.CSRF.TrustedOrigins),
	)

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
