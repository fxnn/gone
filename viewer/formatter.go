package viewer

import (
	"io"
	"mime"
	"net/http"

	"github.com/fxnn/gone/store"
)

type formatter interface {
	serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request)
}

var formatterByMimeType = map[string]formatter{
	store.MarkdownMimeType: newMarkdownFormatter(),
}

func mimeTypeFormatter(mediaType string) formatter {
	if mimeType, _, err := mime.ParseMediaType(mediaType); err == nil {
		if f, ok := formatterByMimeType[mimeType]; ok {
			return f
		}
	}
	return newRawFormatter(mediaType)
}
