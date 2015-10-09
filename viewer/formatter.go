package viewer

import (
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/fxnn/gone/filer"
)

type formatter interface {
	serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request)
}

var formatterByMimeType = map[string]formatter{
	filer.MarkdownMimeType: newMarkdownFormatter(),
}

func mimeTypeFormatter(mediaType string) formatter {
	if mimeType, _, err := mime.ParseMediaType(mediaType); err == nil {
		if f, ok := formatterByMimeType[mimeType]; ok {
			return f
		}
	}
	return newRawFormatter(mediaType)
}
