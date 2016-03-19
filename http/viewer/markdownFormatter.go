package viewer

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/fxnn/gone/log"

	"github.com/fxnn/gone/http/failer"
	"github.com/fxnn/gone/http/templates"
	"github.com/russross/blackfriday"
)

const markdownFormatterOutputMimeType = "text/html"

type markdownFormatter struct {
	renderer *templates.ViewerRenderer
}

func newMarkdownFormatter(l templates.Loader) markdownFormatter {
	// TODO: Preinitialize Markdown Renderer
	var result = markdownFormatter{templates.NewViewerRenderer()}
	if err := result.renderer.LoadAndWatch(l); err != nil {
		panic(fmt.Errorf("couldn't load viewer template: %s", err))
	}
	return result
}

func (f markdownFormatter) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", markdownFormatterOutputMimeType)

	markdown, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Warnf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	html := blackfriday.MarkdownCommon(markdown)
	if err := f.renderer.Render(writer, request.URL, string(html)); err != nil {
		log.Warnf("%s %s: %s", request.Method, request.URL, err)
	}
}
