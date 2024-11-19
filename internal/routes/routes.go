package routes

import (
	"net/http"

	"github.com/Hiwiii/snippetbox.git/config"
	"github.com/Hiwiii/snippetbox.git/internal/handlers"
	"github.com/Hiwiii/snippetbox.git/internal/middleware"
)

// Routes returns an http.ServeMux containing all application routes.
func Routes(app *config.Application, helpers *middleware.Helpers) *http.ServeMux {
	mux := http.NewServeMux()

	// File server for static assets
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register routes using closures
	mux.HandleFunc("/", handlers.Home(app, helpers))
	mux.HandleFunc("/snippet/view", handlers.SnippetView(app, helpers))
	mux.HandleFunc("/snippet/create", handlers.SnippetCreate(app, helpers))

	return mux
}
