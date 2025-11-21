package main

import (
	"context"
	"firecrest-go/ui/templates"
	"firecrest-go/ui/templates/auth"
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	authors, err := app.db.ListAuthors(context.Background())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, http.StatusOK, templates.Home(authors))
}

func (app *application) eventView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	slugId, err := strconv.Atoi(slug)
	if err != nil {

		return
	}

	fmt.Printf("Event slug: %s\n", slug)
	author, err := app.db.GetAuthor(context.Background(), int64(slugId))
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, http.StatusOK, templates.Event(author))
}

/*
* AUTH HANDLERS
=================
*/
func (app *application) signInView(w http.ResponseWriter, r *http.Request) {

	app.render(w, http.StatusOK, auth.SignIn())
}

func (app *application) signUpView(w http.ResponseWriter, r *http.Request) {

	app.render(w, http.StatusOK, auth.SignUp())
}
