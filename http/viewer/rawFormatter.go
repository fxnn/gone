package viewer

import (
	"io"
	"net/http"

	"github.com/fxnn/gone/log"

	"github.com/fxnn/gone/http/failer"
)

type rawFormatter struct {
	mimeType string
}

func newRawFormatter(mimeType string) rawFormatter {
	return rawFormatter{mimeType}
}

func (f rawFormatter) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", f.mimeType)

	// TODO: Use http.ServeContent instead
	_, err := io.Copy(writer, reader)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}
}
