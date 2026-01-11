package main

import (
	"context"
	"errors"
	"firecrest-go/tutorial"
	"firecrest-go/ui/templates"
	"firecrest-go/ui/templates/auth"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	events, err := app.db.ListEvents(context.Background())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, http.StatusOK, templates.Home(events))
}

func (app *application) eventView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	if len(slug) == 0 || len(slug) > 100 {
		app.serverError(w, r, errors.New("invalid event slug"))
		return
	}

	event, err := app.db.GetEvent(context.Background(), slug)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, http.StatusOK, templates.Event(event))
}

/*
* AUTH HANDLERS
=================
*/
func (app *application) signInView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, auth.SignIn())
}

func (app *application) signInPost(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, auth.SignIn())
}

func (app *application) signUpView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, auth.SignUp())
}

func (app *application) signUpPost(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, auth.SignUp())
}

func (app *application) adminCreateView(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, auth.SignIn())
}

func (app *application) adminCreatePost(w http.ResponseWriter, r *http.Request) {
	event, err := app.db.CreateEvent(context.Background(), tutorial.CreateEventParams{
		OrganisationID: 1,
		Name:           "Lincoln 10k",
		Slug:           "lincoln-10k",
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, http.StatusOK, templates.Event(event))
}

// func (app *application) adminCreateOrg(w http.ResponseWriter, r *http.Request) {
// 	org, err := app.db.CreateOrganisation(context.Background(), tutorial.cr{
// 		Name: "Lincoln 10k",
// 		Slug: "lincoln-10k",
// 	})
// 	if err != nil {
// 		app.serverError(w, r, err)
// 		return
// 	}

// 	app.render(w, http.StatusOK, templates.Event(event))
// }

func (app *application) adminCreateUser(w http.ResponseWriter, r *http.Request) {
	_, err := app.db.CreateUser(context.Background(), tutorial.CreateUserParams{
		Email:     "user@example.com",
		FirstName: "Kristian",
		LastName:  "Roebuck",
		Role:      "admin",
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, http.StatusOK, templates.Home(make([]tutorial.Event, 0)))
}
