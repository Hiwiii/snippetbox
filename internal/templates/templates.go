package templates
import (
	"html/template"
	"path/filepath"
	"time"

    "github.com/Hiwiii/snippetbox.git/internal/models"
)

// humanDate formats a time.Time object into a human-readable format.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// functions is a global template.FuncMap object where we register custom functions.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// TemplateData holds the dynamic data passed to HTML templates.
type TemplateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
}

// NewTemplateCache initializes and returns a map of cached templates.
func NewTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Get a list of all page templates in the `./ui/html/pages` folder.
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through each file path.
	for _, page := range pages {
		// Extract the file name (e.g., "home.tmpl").
		name := filepath.Base(page)

		// Create a slice with the base, navigation partial, and current page template.
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}

		// Parse the template files and attach the custom functions using the Funcs method.
		ts, err := template.New(name).Funcs(functions).ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// Add the template to the cache with its name as the key.
		cache[name] = ts
	}

	// Return the cache map.
	return cache, nil
}
