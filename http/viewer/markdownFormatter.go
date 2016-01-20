package viewer

import (
	"io"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	"github.com/fxnn/gone/log"

	"github.com/fxnn/gone/http/failer"
	"github.com/fxnn/gone/http/templates"
	"github.com/russross/blackfriday"
)

const markdownFormatterOutputMimeType = "text/html"

type markdownFormatter struct {
	templateValue *atomic.Value // contain a template.ViewerTemplate
}

func newMarkdownFormatter(l templates.Loader) markdownFormatter {
	// TODO: Preinitialize Markdown Renderer
	var result = markdownFormatter{new(atomic.Value)}
	result.templateValue.Store(templates.LoadViewerTemplate(l))
	go result.watchTemplate(l)
	return result
}

func (f markdownFormatter) watchTemplate(l templates.Loader) {
	for newTemplate := range templates.WatchViewerTemplate(l) {
		f.templateValue.Store(newTemplate)
	}
}

func (f markdownFormatter) template() *templates.ViewerTemplate {
	if result, ok := f.templateValue.Load().(templates.ViewerTemplate); ok {
		return &result
	}
	panic("markdownFormatter template value is of wrong type")
}

func (f markdownFormatter) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", markdownFormatterOutputMimeType)

	markdown, err := ioutil.ReadAll(reader)
	if err == nil {
		html := blackfriday.MarkdownCommon(markdown)
		f.template().Render(writer, request.URL, string(html))
	} else {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
	}
}
