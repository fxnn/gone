package editor

import (
	"fmt"
	"github.com/fxnn/gone/resources"
	"html/template"
	"io"
)

const editTemplateName string = "/edit.html"
const useLocalEditTemplate bool = false

type editTemplate struct {
	template *template.Template
}

func loadEditTemplate() (editTemplate, error) {
	templateString, err := resources.FSString(useLocalEditTemplate, editTemplateName)
	if err != nil {
		return editTemplate{nil}, fmt.Errorf("couldnt load template %s: %s", editTemplateName, err)
	}

	template, err := template.New(editTemplateName).Parse(templateString)
	if err != nil {
		return editTemplate{nil}, fmt.Errorf("couldnt parse template %s: %s", editTemplateName, err)
	}
	return editTemplate{template}, nil
}

func (t *editTemplate) Execute(writer io.Writer, data interface{}) error {
	return t.template.Execute(writer, data)
}
