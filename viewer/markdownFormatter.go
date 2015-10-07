package viewer

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fxnn/gone/failer"
	"github.com/fxnn/gone/internal/github.com/russross/blackfriday"
)

const markdownFormatterOutputMimeType = "text/html"

type markdownFormatter struct{}

func newMarkdownFormatter() markdownFormatter {
	// TODO: Preinitialize Markdown Renderer
	return markdownFormatter{}
}

func (f markdownFormatter) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", markdownFormatterOutputMimeType)

	var (
		markdown []byte
		written  int
		err      error
	)

	markdown, err = ioutil.ReadAll(reader)
	if err == nil {
		written, err = writer.Write(blackfriday.MarkdownCommon(markdown))
	}
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}
	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, written)
}
