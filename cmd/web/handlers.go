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
	app.render(r.Context(), w, http.StatusOK, auth.SignIn())
}

func (app *application) signInPost(w http.ResponseWriter, r *http.Request) {
	app.render(r.Context(), w, http.StatusOK, auth.SignIn())
}

func (app *application) signUpView(w http.ResponseWriter, r *http.Request) {
	app.render(r.Context(), w, http.StatusOK, auth.SignUp())
}

func (app *application) signUpPost(w http.ResponseWriter, r *http.Request) {
	app.render(r.Context(), w, http.StatusOK, auth.SignUp())
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
