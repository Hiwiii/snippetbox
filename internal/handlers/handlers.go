package handlers

import (
	"errors"
	"fmt"
	// "html/template"
	"net/http"
	"strconv"

	"github.com/Hiwiii/snippetbox.git/config"
	"github.com/Hiwiii/snippetbox.git/internal/middleware"
	"github.com/Hiwiii/snippetbox.git/internal/models"
)

// Home handler with dependency injection using middleware.Helpers
func Home(app *config.Application, helpers *middleware.Helpers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the requested path is exactly "/".
		if r.URL.Path != "/" {
			helpers.NotFound(w)
			return
		}

		// Fetch the latest snippets using the SnippetModel's Latest() method.
		snippets, err := app.SnippetModel.Latest()
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		// Print the snippets for testing purposes
        for _, snippet := range snippets {
            fmt.Fprintf(w, "%+v\n", snippet)
        }

		// Define the files for the template.
		// files := []string{
		// 	"./ui/html/base.tmpl",
		// 	"./ui/html/partials/nav.tmpl",
		// 	"./ui/html/pages/home.tmpl",
		// }

		// // Parse the template files.
		// ts, err := template.ParseFiles(files...)
		// if err != nil {
		// 	helpers.ServerError(w, err)
		// 	return
		// }

		// // Pass the snippets data to the template for rendering.
		// err = ts.ExecuteTemplate(w, "base", snippets)
		// if err != nil {
		// 	helpers.ServerError(w, err)
		// }
	}
}


// SnippetCreate handler with dependency injection using middleware.Helpers
func SnippetCreate(app *config.Application, helpers *middleware.Helpers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			helpers.ClientError(w, http.StatusMethodNotAllowed)
			return
		}

		// Create some variables holding dummy data
		// These will be removed during the build in favor of real form data
		title := "O snail"
		content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
		expires := 7

		// Pass the data to the SnippetModel.Insert() method
		id, err := app.SnippetModel.Insert(title, content, expires)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		// Redirect the user to the relevant page for the snippet
		http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
	}
}

// SnippetView handler with dependency injection using middleware.Helpers
// SnippetView handler with dependency injection
func SnippetView(app *config.Application, helpers *middleware.Helpers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the "id" parameter from the URL query string.
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil || id < 1 {
			// If the ID is invalid or less than 1, return a 404 Not Found response.
			helpers.NotFound(w)
			return
		}

		// Use the SnippetModel's Get method to fetch the snippet data.
		snippet, err := app.SnippetModel.Get(id)
		if err != nil {
			if errors.Is(err, models.ErrNoRecord) {
				// If no matching record is found, return a 404 Not Found response.
				helpers.NotFound(w)
			} else {
				// For other errors, return a 500 Internal Server Error response.
				helpers.ServerError(w, err)
			}
			return
		}

		// Write the snippet data as a plain-text HTTP response.
		// This is temporary; you can later modify this to render a proper HTML page.
		fmt.Fprintf(w, "%+v", snippet)
	}
}
