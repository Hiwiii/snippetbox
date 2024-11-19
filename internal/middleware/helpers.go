package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"log"
)

type Helpers struct {
	ErrorLog *log.Logger
}

// serverError writes an error message and stack trace to the error log,
// then sends a generic 500 Internal Server Error response to the user.
func (h *Helpers) ServerError(w http.ResponseWriter, err error) {
    // Format the error message with the stack trace.
    trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
    
    // Use the Output() method to include the file name and line number of the caller.
    h.ErrorLog.Output(2, trace)

    // Send the generic 500 Internal Server Error response to the client.
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}


// clientError sends a specific status code and corresponding description to the user.
// For example, we'll use this later in the book to send responses like 400 "Bad Request"
// when there's a problem with the request that the user sent.
func (h *Helpers) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound is a convenience wrapper around clientError which sends a 404 Not Found response to the user.
func (h *Helpers) NotFound(w http.ResponseWriter) {
	h.ClientError(w, http.StatusNotFound)
}
