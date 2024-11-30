package main

import (
	"crypto/tls" // Import for TLS configuration
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Hiwiii/snippetbox.git/config"
	"github.com/Hiwiii/snippetbox.git/internal/middleware"
	"github.com/Hiwiii/snippetbox.git/internal/models"
	"github.com/Hiwiii/snippetbox.git/internal/routes"
	"github.com/Hiwiii/snippetbox.git/internal/templates"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
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

	// Initialize a new template cache
	templateCache, err := templates.NewTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true // Ensure cookies are only sent over HTTPS

	// Initialize the Application struct
	app := &config.Application{
		InfoLog:        infoLog,
		ErrorLog:       errorLog,
		DB:             db,
		SnippetModel:   snippetModel,
		TemplateCache:  templateCache,
		FormDecoder:    form.NewDecoder(),
		SessionManager: sessionManager,
	}

	// Initialize the Helpers struct
	helpers := &middleware.Helpers{
		ErrorLog:       errorLog,
		TemplateCache:  templateCache,
		FormDecoder:    form.NewDecoder(),
		SessionManager: sessionManager,
	}

	// Configure the TLS settings
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}, // Use efficient elliptic curves
	}

	// Initialize and start the HTTPS server
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   routes.Routes(app, helpers), // Use the Routes function from the routes package
		TLSConfig: tlsConfig,  // Apply the TLS configuration

		// Add Idle, Read, and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem") // Use the TLS certificate and key
	if err != nil {
		errorLog.Fatalf("Could not start server: %v", err)
	}
}
