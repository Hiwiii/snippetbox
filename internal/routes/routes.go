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

	// Create a dynamic middleware chain.
	dynamic := alice.New(helpers.SessionManager.LoadAndSave)

	// Register dynamic routes (routes needing middleware for session handling).
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(handlers.Home(app, helpers)))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(handlers.SnippetView(app, helpers)))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(handlers.SnippetCreate(app, helpers)))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(handlers.SnippetCreatePost(app, helpers)))

	// Create a standard middleware chain for logging, recovery, and headers.
	standard := alice.New(
		func(h http.Handler) http.Handler {
			return middleware.RecoverPanic(app, helpers, h)
		},
		middleware.LogRequest(app),
		middleware.SecureHeaders,
	)

	// Wrap the router with standard middleware and return.
	return standard.Then(router)
}
