package middleware

import (
	"fmt"
	"net/http"

	"github.com/Hiwiii/snippetbox.git/config"
)

// secureHeaders is a middleware that sets various security-related headers.
func SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Content-Security-Policy header to restrict resources like scripts and styles.
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		// Set Referrer-Policy header to control the referrer information sent with requests.
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		// Set X-Content-Type-Options header to prevent MIME type sniffing.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Set X-Frame-Options header to prevent clickjacking attacks.
		w.Header().Set("X-Frame-Options", "deny")

		// Set X-XSS-Protection header to disable the deprecated browser XSS filter.
		w.Header().Set("X-XSS-Protection", "0")

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// LogRequest logs information about each HTTP request.
func LogRequest(app *config.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Log the IP address, protocol, HTTP method, and requested URL.
			app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

			// Call the next handler in the chain.
			next.ServeHTTP(w, r)
		})
	}
}

// RecoverPanic recovers from panics and returns a 500 Internal Server Error.
func RecoverPanic(app *config.Application, helpers *Helpers, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")

				// Log the error and return a 500 error response.
				app.ErrorLog.Printf("Panic: %v", err)
				app.ErrorLog.Output(2, fmt.Sprintf("%s", err))
				helpers.ServerError(w, fmt.Errorf("%s", err))
			}
		}()

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
