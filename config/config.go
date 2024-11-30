package config

import (
	"database/sql"
	"github.com/Hiwiii/snippetbox.git/internal/models"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
)

// Application holds the dependencies for the application.
type Application struct {
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	DB             *sql.DB
	SnippetModel   *models.SnippetModel
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}
