package templates

import (
	"fmt"
	"io"
	"net/url"
)

const editorTemplateName string = "/editor.html"

// EditorRenderer renders the edit UI from a go HTML template.
// Always check the Err() result!
type EditorRenderer struct {
	*renderer
}

func NewEditorRenderer() *EditorRenderer {
	return &EditorRenderer{
		newRenderer(editorTemplateName),
	}
}

func (r EditorRenderer) Render(writer io.Writer, url *url.URL, content string,
	mimeType string, edit bool) error {
	var data = make(map[string]interface{})
	data["path"] = url.Path
	data["content"] = content
	data["contenttype"] = mimeType
	if edit {
		data["edit"] = "edit"
	}

	if err := r.renderData(writer, data); err != nil {
		return fmt.Errorf("couldn't render editor template: %s", err)
	}

	return nil
}
