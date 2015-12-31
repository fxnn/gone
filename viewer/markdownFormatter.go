package viewer

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fxnn/gone/failer"
	"github.com/fxnn/gone/templates"
	"github.com/russross/blackfriday"
)

const markdownFormatterOutputMimeType = "text/html"

type markdownFormatter struct {
	template templates.ViewerTemplate
}

func newMarkdownFormatter(l templates.Loader) markdownFormatter {
	// TODO: Preinitialize Markdown Renderer
	return markdownFormatter{templates.LoadViewerTemplate(l)}
}

func (f markdownFormatter) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", markdownFormatterOutputMimeType)

	markdown, err := ioutil.ReadAll(reader)
	if err == nil {
		html := blackfriday.MarkdownCommon(markdown)
		f.template.Render(writer, request.URL, string(html))
		log.Printf("%s %s: delivered markdown page", request.Method, request.URL)
	} else {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
	}
}
