package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Hiwiii/snippetbox.git/config"
	"github.com/Hiwiii/snippetbox.git/internal/middleware"
	"github.com/Hiwiii/snippetbox.git/internal/routes"
	"github.com/Hiwiii/snippetbox.git/internal/models"
)

func main() {
	// Define flags for the server address and DSN (data source name)
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:Secure@123@tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL DSN")
	flag.Parse()

	// Create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Open the database connection
	db, err := config.OpenDB(*dsn)
	if err != nil {
		errorLog.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close() // Ensure the connection is closed when the program exits

	// Initialize the SnippetModel with the database connection
	snippetModel := &models.SnippetModel{DB: db}
	// Initialize the Application struct
	app := &config.Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		DB:       db,
		SnippetModel: snippetModel,
	}

	// Initialize the Helpers struct
	helpers := &middleware.Helpers{
		ErrorLog: errorLog,
	}

	// Initialize and start the HTTP server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  routes.Routes(app, helpers), // Use the Routes function from the routes package
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
	
}
