package templates

// Loader encapsulates how to retrieve a template by its name.
// Concrete implementations could load data packaged with the binary, files
// from filesystem or even network resources.
// Caching is not the duty of a Loader.
type Loader interface {
	// LoadHtmlTemplate loads the template with the given name.
	LoadHtmlTemplate(name string) Template

	// WatchHtmlTemplate returns a chan over which changed Templates might be
	// received, if supported by this loader.
	WatchHtmlTemplate(name string) <-chan Template

	// Close frees resources bound by this loader, esp. in case of watching.
	Close() error
}
