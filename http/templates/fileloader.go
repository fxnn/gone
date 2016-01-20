package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/fxnn/gopath"
)

// FilesystemLoader is a Provider that loads templates from the filesystem.
// It only loads the template once and then holds it in memory.
// It reloads the template after the file changed.
type FilesystemLoader struct {
	root gopath.GoPath
}

// NewFilesystemLoader creates a new instance with templates located in the
// given root path.
func NewFilesystemLoader(root gopath.GoPath) *FilesystemLoader {
	if root.HasErr() {
		panic(fmt.Sprintf("NewFilesystemLoader: root has error: %s", root.Err()))
	}
	if !root.IsExists() {
		panic(fmt.Sprintf("NewFilesystemLoader: root %s does not exist", root.Path()))
	}

	return &FilesystemLoader{root}
}

func (l *FilesystemLoader) templatePath(name string) gopath.GoPath {
	return l.root.JoinPath(name)
}

func (l *FilesystemLoader) LoadHtmlTemplate(name string) Template {
	p := l.templatePath(name)
	contentBytes, err := ioutil.ReadFile(p.Path())
	if err != nil {
		return newWithError(fmt.Errorf("couldn't load template %s: %s", p.Path(), err))
	}

	htmlTemplate, err := template.New(name).Parse(string(contentBytes))
	if err != nil {
		return newWithError(fmt.Errorf("couldn't parse template %s: %s", p.Path(), err))
	}
	return newFromHtmlTemplate(htmlTemplate)
}
