package viewer

import (
	"io"
	"mime"
	"net/http"

	"github.com/fxnn/gone/http/templates"
	"github.com/fxnn/gone/store"
)

type formatter interface {
	serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request)
}

type formatters struct {
	formatterByMimeType map[string]formatter
}

func newFormatters(l templates.Loader) formatters {
	var formatterByMimeType = map[string]formatter{
		store.MarkdownMimeType: newMarkdownFormatter(l),
	}
	return formatters{formatterByMimeType}
}

func (s *formatters) mimeTypeFormatter(mediaType string) formatter {
	if mimeType, _, err := mime.ParseMediaType(mediaType); err == nil {
		if f, ok := s.formatterByMimeType[mimeType]; ok {
			return f
		}
	}
	return newRawFormatter(mediaType)
}
