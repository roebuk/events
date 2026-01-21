package main

import (
	"errors"
	"net/http"

	"firecrest/internal/repository"
	"firecrest/internal/service"
	"firecrest/ui/templates"
	"firecrest/ui/templates/auth"
)

func (app *application) health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		app.logger.Error("failed to write health response", "error", err)
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	events, err := app.eventService.ListEvents(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(r.Context(), w, http.StatusOK, templates.Home(events))
}

func (app *application) eventView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	event, err := app.eventService.GetEvent(r.Context(), slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			app.notFound(w)
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		app.serverError(w, r, err)
		return
	}

	app.render(r.Context(), w, http.StatusOK, templates.Event(event))
}

/*
* AUTH HANDLERS
=================
*/
func (app *application) signInView(w http.ResponseWriter, r *http.Request) {
	flashes := app.getAllFlashes(r)
	app.render(r.Context(), w, http.StatusOK, auth.SignIn(flashes))
}

func (app *application) signInPost(w http.ResponseWriter, r *http.Request) {
	// Parse form
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	rememberMe := r.PostForm.Get("remember_me") == "on"

	// Authenticate user
	result, err := app.authService.SignIn(r.Context(), service.SignInInput{
		Email:      email,
		Password:   password,
		RememberMe: rememberMe,
	})

	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			app.addFlash(r, FlashError, "Invalid email or password")
		case errors.Is(err, service.ErrEmailNotVerified):
			app.addFlash(r, FlashWarning, "Please verify your email address before signing in")
		case errors.Is(err, service.ErrAccountLocked):
			app.addFlash(r, FlashError, "Your account has been locked due to too many failed login attempts. Please try again later.")
		case errors.Is(err, service.ErrInvalidInput):
			app.addFlash(r, FlashError, "Please provide both email and password")
		default:
			app.serverError(w, r, err)
			return
		}
		http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
		return
	}

	// Regenerate session token
	if err := app.sessionManager.RenewToken(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	// Store userID in session
	app.sessionManager.Put(r.Context(), "userID", result.User.ID)

	app.addFlash(r, FlashSuccess, "Welcome back!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signUpView(w http.ResponseWriter, r *http.Request) {
	flashes := app.getAllFlashes(r)
	app.render(r.Context(), w, http.StatusOK, auth.SignUp(flashes))
}

func (app *application) signUpPost(w http.ResponseWriter, r *http.Request) {
	// Parse form
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	firstName := r.PostForm.Get("first_name")
	lastName := r.PostForm.Get("last_name")

	// Create user
	_, err := app.authService.SignUp(r.Context(), service.SignUpInput{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	})

	if err != nil {
		// Handle specific errors
		switch {
		case errors.Is(err, service.ErrEmailExists):
			app.addFlash(r, FlashError, "An account with this email already exists")
		case errors.Is(err, service.ErrInvalidInput):
			app.addFlash(r, FlashError, err.Error())
		default:
			app.serverError(w, r, err)
			return
		}
		http.Redirect(w, r, "/auth/sign-up", http.StatusSeeOther)
		return
	}

	app.addFlash(r, FlashSuccess, "Account created successfully! Please check your email to verify your account.")
	http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
}

func (app *application) signOut(w http.ResponseWriter, r *http.Request) {
	// Destroy the session
	if err := app.sessionManager.Destroy(r.Context()); err != nil {
		app.serverError(w, r, err)
		return
	}

	app.addFlash(r, FlashInfo, "You have been signed out")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) adminCreatePost(w http.ResponseWriter, r *http.Request) {
	event, err := app.eventService.CreateEvent(r.Context(), service.CreateEventInput{
		OrganisationID: 1,
		Name:           "Lincoln 10k",
		Slug:           "lincoln-10k",
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(r.Context(), w, http.StatusOK, templates.Event(event))
}

func (app *application) adminCreateUser(w http.ResponseWriter, r *http.Request) {
	_, err := app.userService.CreateUser(r.Context(), service.CreateUserInput{
		Email:     "user@example.com",
		FirstName: "Kristian",
		LastName:  "Roebuck",
		Role:      "admin",
	})

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(r.Context(), w, http.StatusOK, templates.Home(nil))
}
