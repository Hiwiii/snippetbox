package routes

import (
	"net/http"

	"github.com/Hiwiii/snippetbox.git/config"
	"github.com/Hiwiii/snippetbox.git/internal/handlers"
	"github.com/Hiwiii/snippetbox.git/internal/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Routes returns an http.Handler with all application routes wrapped in middleware.
func Routes(app *config.Application, helpers *middleware.Helpers) http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// Set a custom NotFound handler to use the helpers.NotFound function.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		helpers.NotFound(w)
	})

	// File server for static files.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Register routes with handlers using dynamic parameters.
	router.HandlerFunc(http.MethodGet, "/", handlers.Home(app, helpers))
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", handlers.SnippetView(app, helpers))
    router.HandlerFunc(http.MethodGet, "/snippet/create", handlers.SnippetCreate(app, helpers)) // For displaying the form
    router.HandlerFunc(http.MethodPost, "/snippet/create", handlers.SnippetCreatePost(app, helpers)) // For handling form submission


	// Create the middleware chain using alice.
	standard := alice.New(
		func(h http.Handler) http.Handler {
			return middleware.RecoverPanic(app, helpers, h)
		},
		middleware.LogRequest(app),
		middleware.SecureHeaders,
	)

	// Wrap the router with middleware and return.
	return standard.Then(router)
}
