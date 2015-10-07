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

func LoadViewerTemplate() ViewerTemplate {
	return ViewerTemplate{
		loadHtmlTemplate(viewerTemplateName),
	}
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
