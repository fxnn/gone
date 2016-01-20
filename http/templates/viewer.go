package templates

import (
	"html/template"
	"io"
	"net/url"
)

const viewerTemplateName string = "/viewer.html"

// Always check the Err() result!
type ViewerTemplate struct {
	Template
}

func LoadViewerTemplate(loader Loader) ViewerTemplate {
	return NewViewerTemplate(
		loader.LoadHtmlTemplate(viewerTemplateName),
	)
}

func NewViewerTemplate(t Template) ViewerTemplate {
	return ViewerTemplate{t}
}

func WatchViewerTemplate(loader Loader) <-chan ViewerTemplate {
	var viewerTemplateChan = make(chan ViewerTemplate)
	go func() {
		for template := range loader.WatchHtmlTemplate(viewerTemplateName) {
			viewerTemplateChan <- NewViewerTemplate(template)
		}
		close(viewerTemplateChan)
	}()
	return viewerTemplateChan
}

func (t *ViewerTemplate) Render(writer io.Writer, url *url.URL, htmlContent string) {
	if t.err != nil {
		return
	}

	var data = make(map[string]interface{})
	data["path"] = url.Path
	data["htmlContent"] = template.HTML(htmlContent)

	t.Execute(writer, data)
}
