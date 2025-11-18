package main

import (
	"context"
	"firecrest-go/ui/templates"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	authors, err := app.db.ListAuthors(context.Background())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, http.StatusOK, templates.Authors(authors))
	// }

	// app.render(w, r, http.StatusOK, templates.Authors(authors))
}

func (app *application) eventView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	w.Write([]byte("Event details page for slug: " + slug))
}
