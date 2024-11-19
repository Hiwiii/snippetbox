package config

import (
	"database/sql"
	"log"
	"github.com/Hiwiii/snippetbox.git/internal/models"
)

// Application holds the dependencies for the application.
type Application struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DB       *sql.DB
	SnippetModel *models.SnippetModel
}