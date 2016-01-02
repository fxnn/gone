package templates

// Loader encapsulates how to load a template by its name.
// Concrete implementations could load data packaged with the binary, files
// from filesystem or even network resources.
type Loader interface {
	LoadHtmlTemplate(name string) Template
}
