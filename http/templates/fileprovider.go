package templates

import (
	"sync/atomic"

	"gopkg.in/fsnotify.v1"
)

type FilesystemProvider struct {
	name    string            // Name of the template to load
	loader  Loader            // Template loader
	watcher *fsnotify.Watcher // Notifies about changed files
	value   atomic.Value      // Currently loaded template
}

// NewFilesystemProvider constructs a provider for the given template and the
// given FilesystemLoader.
func NewFilesystemProvider(name string, loader *FilesystemLoader) (*FilesystemProvider, error) {
	var watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(loader.templatePath(name).Path())
	if err != nil {
		return nil, err
	}

	var value atomic.Value
	value.Store(loader.LoadHtmlTemplate(name))

	return &FilesystemProvider{name, loader, watcher, value}, nil
}

func (p *FilesystemProvider) Close() error {
	return p.watcher.Close()
}

func (p *FilesystemProvider) HtmlTemplate() Template {
	return p.value.Load().(Template)
}
