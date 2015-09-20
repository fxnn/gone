package templates

import (
	"fmt"
	"github.com/fxnn/gone/resources"
	"html/template"
	"io"
)

const useLocalEditTemplate bool = false

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
func (t Template) Err() error {
	return t.err
}

func loadHtmlTemplate(name string) Template {
	content, err := resources.FSString(useLocalEditTemplate, name)
	if err != nil {
		return newWithError(fmt.Errorf("couldnt load template %s: %s", name, err))
	}

	htmlTemplate, err := template.New(name).Parse(content)
	if err != nil {
		return newWithError(fmt.Errorf("couldnt parse template %s: %s", name, err))
	}
	return newFromHtmlTemplate(htmlTemplate)
}

func (t *Template) Execute(writer io.Writer, data interface{}) {
	if t.err != nil {
		return
	}
	t.err = t.template.Execute(writer, data)
}
