package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Hiwiii/snippetbox.git/config"
	"github.com/Hiwiii/snippetbox.git/internal/forms"
	"github.com/Hiwiii/snippetbox.git/internal/middleware"
	"github.com/Hiwiii/snippetbox.git/internal/models"
	"github.com/Hiwiii/snippetbox.git/internal/validators"

	"github.com/julienschmidt/httprouter"
	// "github.com/Hiwiii/snippetbox.git/internal/templates"
)

// Home handler with dependency injection using middleware.Helpers
func Home(app *config.Application, helpers *middleware.Helpers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snippets, err := app.SnippetModel.Latest()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data := helpers.NewTemplateData(r)
		data.Snippets = snippets

		helpers.Render(w, http.StatusOK, "home.tmpl", data)
	}
}

// SnippetCreateForm handler with dependency injection using middleware.Helpers
func SnippetCreate(app *config.Application, helpers *middleware.Helpers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Initialize template data
		data := helpers.NewTemplateData(r)

		// Set the default snippet expiry to 365 days.
		data.Form = forms.SnippetCreateForm{
			Expires: 365,
		}

		// Render the form template
		helpers.Render(w, http.StatusOK, "create.tmpl", data)
	}
}

// SnippetCreatePost handler with dependency injection using middleware.Helpers
func SnippetCreatePost(app *config.Application, helpers *middleware.Helpers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Declare a new empty instance of the SnippetCreateForm struct.
		var form forms.SnippetCreateForm

		// Use the DecodePostForm helper to decode the form data into the struct.
		err := helpers.DecodePostForm(r, &form)
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		// Validate the form fields using the validator.
		form.Validator.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
		form.Validator.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")
		form.Validator.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
		form.Validator.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

		// If validation fails, re-display the form with validation errors.
		if !form.Validator.Valid() {
			data := helpers.NewTemplateData(r)
			data.Form = form
			helpers.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		// Pass the validated form data to the SnippetModel.Insert() method.
		id, err := app.SnippetModel.Insert(form.Title, form.Content, form.Expires)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		// Use the SessionManager to add a flash message to the session.
		app.SessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

		// Redirect the user to the relevant page for the snippet.
		http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
	}
}

// SnippetView handler with dependency injection using middleware.Helpers
func SnippetView(app *config.Application, helpers *middleware.Helpers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil || id < 1 {
			helpers.NotFound(w)
			return
		}

		snippet, err := app.SnippetModel.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				helpers.NotFound(w)
			} else {
				helpers.ServerError(w, err)
			}
			return
		}

		data := helpers.NewTemplateData(r)
		data.Snippet = snippet

		helpers.Render(w, http.StatusOK, "view.tmpl", data)
	}
}
