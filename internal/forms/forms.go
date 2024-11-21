package forms

type SnippetCreateForm struct {
    Title       string
    Content     string
    Expires     int
    FieldErrors map[string]string
}
