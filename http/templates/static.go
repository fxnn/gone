package templates

import (
	"fmt"
	"html/template"
	"os"

	"github.com/fxnn/gone/resources"
	"github.com/fxnn/gopath"
)

// StaticLoader is a Provider that loads templates from data packaged with the
// application binary.
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

func (l *StaticLoader) WriteAllTemplates(targetDir gopath.GoPath) error {
	if err := os.MkdirAll(targetDir.Path(), 0777); err != nil {
		return fmt.Errorf("couldn't create dir %s: %s", targetDir, err)
	}

	for _, name := range []string{"/editor.html", "/viewer.html"} {
		var targetFile = targetDir.JoinPath(name)
		if targetFile.HasErr() {
			return fmt.Errorf("couldn't create path for template %s: %s", name, targetFile.Err())
		}

		content, err := resources.FSString(l.useLocalTemplates, name)
		if err != nil {
			return fmt.Errorf("couldn't open template %s: %s", name, err)
		}

		out, err := os.Create(targetFile.Path())
		if err != nil {
			return fmt.Errorf("couldn't create file %s: %s", targetFile, err)
		}

		out.WriteString(content)
		if out.Close(); err != nil {
			return fmt.Errorf("couldn't close file %s: %s", targetFile, err)
		}
	}

	return nil
}
