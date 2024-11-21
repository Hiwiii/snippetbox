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
		// Parse the form data.
		err := r.ParseForm()
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		// Get the expires value from the form.
		expires, err := strconv.Atoi(r.PostForm.Get("expires"))
		if err != nil {
			helpers.ClientError(w, http.StatusBadRequest)
			return
		}

		// Initialize a new Validator instance.
		input := validator.Validator{}

		// Perform validation checks using the validator package.
		input.CheckField(validator.NotBlank(r.PostForm.Get("title")), "title", "This field cannot be blank")
		input.CheckField(validator.MaxChars(r.PostForm.Get("title"), 100), "title", "This field cannot be more than 100 characters long")
		input.CheckField(validator.NotBlank(r.PostForm.Get("content")), "content", "This field cannot be blank")
		input.CheckField(validator.PermittedInt(expires, 1, 7, 365), "expires", "This field must equal 1, 7, or 365")

		// If there are validation errors, re-display the form with the errors.
		if !input.Valid() {
			data := helpers.NewTemplateData(r)
			data.Form = struct {
				Title       string
				Content     string
				Expires     int
				FieldErrors map[string]string
			}{
				Title:       r.PostForm.Get("title"),
				Content:     r.PostForm.Get("content"),
				Expires:     expires,
				FieldErrors: input.FieldErrors,
			}
			helpers.Render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
			return
		}

		// Pass the data to the SnippetModel.Insert() method to save it in the database.
		id, err := app.SnippetModel.Insert(r.PostForm.Get("title"), r.PostForm.Get("content"), expires)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

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

