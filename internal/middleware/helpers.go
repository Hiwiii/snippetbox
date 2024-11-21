package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"html/template"
	"log"
	"time"
	"bytes"

	"github.com/Hiwiii/snippetbox.git/internal/templates"
)

type Helpers struct {
	ErrorLog      *log.Logger
	TemplateCache map[string]*template.Template
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

// Render retrieves the appropriate template from the cache and renders it.
// Render retrieves the appropriate template from the cache and renders it.
func (h *Helpers) Render(w http.ResponseWriter, status int, page string, data interface{}) {
	// Retrieve the appropriate template set from the cache.
	ts, ok := h.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		h.ServerError(w, err)
		return
	}

	// Initialize a new buffer to hold the rendered template.
	buf := new(bytes.Buffer)

	// Write the template to the buffer instead of the ResponseWriter.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		h.ServerError(w, err)
		return
	}

	// If successful, write the HTTP status code to the ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the ResponseWriter.
	_, err = buf.WriteTo(w)
	if err != nil {
		h.ServerError(w, err)
	}
}

// NewTemplateData initializes and returns a TemplateData struct.
func (h *Helpers) NewTemplateData(r *http.Request) *templates.TemplateData {
	return &templates.TemplateData{
		CurrentYear: time.Now().Year(),
	}
}
