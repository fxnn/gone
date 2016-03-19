package templates

import (
	"fmt"
	"html/template"
	"io"
	"net/url"
)

const viewerTemplateName string = "/viewer.html"

type ViewerRenderer struct {
	*renderer
}

func NewViewerRenderer() *ViewerRenderer {
	return &ViewerRenderer{newRenderer(viewerTemplateName)}
}

func (r ViewerRenderer) Render(writer io.Writer, url *url.URL, htmlContent string) error {
	var data = make(map[string]interface{})
	data["path"] = url.Path
	data["htmlContent"] = template.HTML(htmlContent)

	if err := r.renderData(writer, data); err != nil {
		return fmt.Errorf("couldn't render viewer template: %s", err)
	}

	return nil
}
