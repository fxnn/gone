package templates

import (
	"fmt"
	"io"
	"net/url"
)

const uploaderTemplateName string = "/uploader.html"

// UploaderRenderer renders the upload UI from a go HTML template.
type UploaderRenderer struct {
	*renderer
}

func NewUploaderRenderer() *UploaderRenderer {
	return &UploaderRenderer{
		newRenderer(uploaderTemplateName),
	}
}

func (r UploaderRenderer) Render(writer io.Writer, url *url.URL,
	targetAlreadyExists bool) error {
	var data = make(map[string]interface{})
	data["path"] = url.Path
	if targetAlreadyExists {
		data["targetAlreadyExists"] = "true"
	}

	if err := r.renderData(writer, data); err != nil {
		return fmt.Errorf("couldn't render uploader template: %s", err)
	}

	return nil
}
