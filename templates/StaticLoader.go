package templates

import (
	"fmt"
	"html/template"

	"github.com/fxnn/gone/resources"
)

// StaticLoader loads templates from data packaged with the application binary.
type StaticLoader struct {
	// useLocalTemplate tells the resource engine to load the templates from the
	// working directory
	useLocalTemplates bool
}

// NewStaticLoader creates a new instance
func NewStaticLoader() *StaticLoader {
	return &StaticLoader{false}
}

// NewStaticLoaderFromWorkingDirectory is to be used for development purposes
// and loads the templates from the application's source directory.
func NewStaticLoaderFromWorkingDirectory() *StaticLoader {
	return &StaticLoader{true}
}

func (l *StaticLoader) LoadHtmlTemplate(name string) Template {
	content, err := resources.FSString(l.useLocalTemplates, name)
	if err != nil {
		return newWithError(fmt.Errorf("couldn't load template %s: %s", name, err))
	}

	htmlTemplate, err := template.New(name).Parse(content)
	if err != nil {
		return newWithError(fmt.Errorf("couldn't parse template %s: %s", name, err))
	}
	return newFromHtmlTemplate(htmlTemplate)
}
