package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"

	"github.com/fxnn/gopath"
)

// FilesystemLoader loads templates from data packaged with the application binary.
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

func (l *FilesystemLoader) LoadHtmlTemplate(name string) Template {
	p := l.root.JoinPath(name)
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
