package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"

	"gopkg.in/fsnotify.v1"

	"github.com/fxnn/gopath"
)

// FilesystemLoader is a Loader that loads templates from the filesystem.
// It only loads the template once and then holds it in memory.
type FilesystemLoader struct {
	root          gopath.GoPath
	watcher       *fsnotify.Watcher
	templateChans map[string]chan Template
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
		make(map[string]chan Template),
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

func (l *FilesystemLoader) WatchHtmlTemplate(name string) <-chan Template {
	var path = l.templatePath(name).Path()
	if err := l.watcher.Add(path); err != nil {
		log.Printf("start watching filesystem template %s: %s", path, err)
		return make(chan Template)
	}

	l.templateNames[path] = name
	templateChan, ok := l.templateChans[path]
	if !ok {
		templateChan = make(chan Template)
		l.templateChans[path] = templateChan
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
				log.Printf("reloading template %s from %s", name, path)
				templateChan <- l.LoadHtmlTemplate(name)
			}
		case err, ok := <-l.watcher.Errors:
			if !ok {
				log.Printf("watching filesystem templates stopped")
				return
			}
			log.Printf("watching filesystem templates: %s", err)
		}
	}
}
