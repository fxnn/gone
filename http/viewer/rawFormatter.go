package viewer

import (
	"io"
	"log"
	"net/http"

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
	written, err := io.Copy(writer, reader)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, written)
}
