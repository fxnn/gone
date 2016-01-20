package templates

import (
	"html/template"
	"io"
)

// Template encapsules basic data and functionality for loading, parsing and
// rendering UI templates.
type Template struct {
	template *template.Template
	err      error
}

func newFromHtmlTemplate(template *template.Template) Template {
	return Template{
		template: template,
		err:      nil,
	}
}

func newWithError(err error) Template {
	return Template{
		template: nil,
		err:      err,
	}
}

// Err() gibt Fehler zurück, die beim Laden oder Rendern des Templates
// auftraten.
// Muss nach jeder Operation ausgefüht werden.
func (t *Template) Err() error {
	var result = t.err
	t.err = nil
	return result
}

func (t *Template) Execute(writer io.Writer, data interface{}) {
	if t.err != nil {
		return
	}
	t.err = t.template.Execute(writer, data)
}
