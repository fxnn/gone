package templates

import (
	"errors"
	"html/template"
	"io"
	"sync/atomic"
)

type renderer struct {
	templateName string
	template     atomic.Value // contains a *template.Template
}

func newRenderer(templateName string) *renderer {
	return &renderer{templateName: templateName}
}

func (r *renderer) Load(l Loader) error {
	if template, err := l.LoadHtmlTemplate(r.templateName); err != nil {
		return err
	} else {
		r.setTemplate(template)
		return nil
	}
}

func (r *renderer) LoadAndWatch(l Loader) error {
	if err := r.Load(l); err != nil {
		return err
	}

	go func() {
		for t := range l.WatchHtmlTemplate(r.templateName) {
			r.setTemplate(t)
		}
	}()

	return nil
}

func (r *renderer) setTemplate(t *template.Template) {
	r.template.Store(t)
}

func (r *renderer) renderData(writer io.Writer, data map[string]interface{}) error {
	if r.template.Load() == nil {
		return errors.New("no template loaded")
	}

	var t = r.template.Load().(*template.Template)
	if err := t.Execute(writer, data); err != nil {
		return err
	}

	return nil
}
