package templates

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/fsnotify.v1"

	"github.com/fxnn/gone/log"
	"github.com/fxnn/gopath"
)

// FilesystemLoader is a Loader that loads templates from the filesystem.
// It supports watching the filesystem for changes in template files.
type FilesystemLoader struct {
	root          gopath.GoPath
	watcher       *fsnotify.Watcher
	templateChans map[string]chan *template.Template
	templateNames map[string]string
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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(fmt.Sprintf("NewFilesystemLoader: can't open watcher: %s", err))
	}

	var loader = &FilesystemLoader{
		root,
		watcher,
		make(map[string]chan *template.Template),
		make(map[string]string)}
	go loader.processEvents()
	return loader
}

func (l *FilesystemLoader) Close() error {
	return l.watcher.Close()
}

func (l *FilesystemLoader) templatePath(name string) gopath.GoPath {
	return l.root.JoinPath(name)
}

func (l *FilesystemLoader) LoadResource(name string) (io.ReadCloser, error) {
	p := l.templatePath(name)
	if p.Err() != nil {
		return nil, fmt.Errorf("couldn't load template resource %s: %s", name, p.Err())
	}

	file, err := os.Open(p.Path())
	if err != nil {
		return nil, fmt.Errorf("couldn't load template resource %s: %s", p.Path(), err)
	}

	return file, nil
}

func (l *FilesystemLoader) LoadHtmlTemplate(name string) (*template.Template, error) {
	p := l.templatePath(name)
	if p.Err() != nil {
		return nil, fmt.Errorf("couldn't load template %s: %s", name, p.Err())
	}

	contentBytes, err := ioutil.ReadFile(p.Path())
	if err != nil {
		return nil, fmt.Errorf("couldn't load template %s: %s", p.Path(), err)
	}

	htmlTemplate, err := template.New(name).Parse(string(contentBytes))
	if err != nil {
		return nil, fmt.Errorf("couldn't parse template %s: %s", p.Path(), err)
	}
	if htmlTemplate == nil {
		return nil, fmt.Errorf("template %s parsed to nil", p.Path())
	}
	return htmlTemplate, nil
}

func (l *FilesystemLoader) WatchHtmlTemplate(name string) <-chan *template.Template {
	var path = l.templatePath(name).Path()

	l.templateNames[path] = name
	templateChan, ok := l.templateChans[path]
	if !ok {
		templateChan = make(chan *template.Template)
		l.templateChans[path] = templateChan
	}

	if err := l.watcher.Add(path); err != nil {
		log.Printf("couldn't watch filesystem template %s: %s", path, err)
		return make(chan *template.Template)
	}

	return templateChan
}

func (l *FilesystemLoader) processEvents() {
	for {
		select {
		case event, ok := <-l.watcher.Events:
			if !ok {
				log.Printf("watching filesystem templates stopped")
				return
			}
			var path, name = event.Name, l.templateNames[event.Name]
			var templateChan = l.templateChans[path]
			if event.Op == fsnotify.Write || event.Op == fsnotify.Chmod {
				if template, err := l.LoadHtmlTemplate(name); err != nil {
					log.Warnf("error while reloading template %s from %s: %s", name, path, err)
				} else {
					log.Printf("reloading template %s from %s", name, path)
					templateChan <- template
				}
			}
		case err, ok := <-l.watcher.Errors:
			if !ok {
				log.Printf("watching filesystem templates stopped")
				return
			}
			log.Printf("error while watching filesystem templates: %s", err)
		}
	}
}
