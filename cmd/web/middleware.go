package main

import (
	"context"
	"net/http"

	"firecrest/db"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"style-src 'self'; "+
				"font-src; "+
				"frame-ancestors 'none'; "+
				"upgrade-insecure-requests")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Permissions-Policy", "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()")
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("request received", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

// requireAuth ensures the user is authenticated.
func (app *application) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			app.addFlash(r, FlashError, "Please sign in to continue")
			http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// redirectIfAuth redirects authenticated users away from auth pages.
func (app *application) redirectIfAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// loadUser loads the authenticated user into the request context.
func (app *application) loadUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.isAuthenticated(r) {
			userID := app.getUserID(r)
			user, err := app.userService.GetUser(r.Context(), userID)
			if err != nil {
				// Session is invalid, clear it
				if err := app.sessionManager.Destroy(r.Context()); err != nil {
					app.logger.Error("failed to destroy session", "error", err)
				}
				next.ServeHTTP(w, r)
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), contextKeyUser, user)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

// Context keys
type contextKey string

const contextKeyUser = contextKey("user")

// getUserFromContext retrieves the user from the request context.
func getUserFromContext(r *http.Request) (db.User, bool) {
	user, ok := r.Context().Value(contextKeyUser).(db.User)
	return user, ok
}
