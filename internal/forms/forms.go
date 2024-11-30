package forms

import "github.com/Hiwiii/snippetbox.git/internal/validators"

type SnippetCreateForm struct {
	Title       string `form:"title"`
	Content     string `form:"content"`
	Expires     int    `form:"expires"`
	FieldErrors map[string]string
	Validator   validator.Validator `form:"-"`
}
