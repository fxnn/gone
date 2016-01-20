package templates

// Loader encapsulates how to retrieve a template by its name.
// Concrete implementations could load data packaged with the binary, files
// from filesystem or even network resources.
// Caching is not the duty of a Loader.
type Loader interface {
	LoadHtmlTemplate(name string) Template
}

// Provider encapsulates how to retrieve an always up-to-date version of a
// specific template.
// It could simply always return the same, once-loaded version, or it could
// reload it once it changes on the file system.
type Provider interface {
	HtmlTemplate() Template
}
