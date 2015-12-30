package templates

import (
	"io"
	"net/url"
)

const editorTemplateName string = "/editor.html"

// Always check the Err() result!
type EditorTemplate struct {
	Template
}

func LoadEditorTemplate() EditorTemplate {
	return EditorTemplate{
		loadHtmlTemplate(editorTemplateName),
	}
}

func (t *EditorTemplate) Render(writer io.Writer, url *url.URL, content string, edit bool) {
	if t.err != nil {
		return
	}

	var data = make(map[string]string)
	data["path"] = url.Path
	data["content"] = content
	if edit {
		data["edit"] = "edit"
	}

	t.Execute(writer, data)
}
