package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/Hiwiii/snippetbox.git/internal/templates"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

type Helpers struct {
	ErrorLog       *log.Logger
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
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
		Flash:       h.SessionManager.PopString(r.Context(), "flash"),
	}
}

// DecodePostForm decodes form data from an HTTP request into a destination struct.
// The second parameter `dst` is the target destination for the decoded data.
func (h *Helpers) DecodePostForm(r *http.Request, dst any) error {
	// Parse the form data to populate r.PostForm.
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Decode the form data into the target destination struct.
	err = h.FormDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// Check if the error is an InvalidDecoderError, which indicates an issue with the target destination.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err) // Panic for InvalidDecoderError to indicate programmer misuse.
		}

		// For all other errors, return them as normal.
		return err
	}

	return nil
}
