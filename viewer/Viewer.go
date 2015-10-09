package viewer

import (
	"github.com/fxnn/gone/failer"
	"github.com/fxnn/gone/filer"
	"log"
	"net/http"
)

const fallbackMimeType = "application/octet-stream"

// The Viewer serves HTTP requests with content from the filesystem.
type Viewer struct {
	filer *filer.Filer
}

// New initializes a Viewer instance ready to use.
func New(filer *filer.Filer) *Viewer {
	return &Viewer{filer}
}

func (v *Viewer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !v.filer.HasReadAccessForRequest(request) {
		log.Printf("%s %s: no read permissions", request.Method, request.URL)
		failer.ServeUnauthorized(writer, request)
		return
	}

	if request.Method != "GET" {
		v.serveNonGET(writer, request)
		return
	}

	v.serveGET(writer, request)
}

func (v *Viewer) serveNonGET(writer http.ResponseWriter, request *http.Request) {
	failer.ServeMethodNotAllowed(writer, request)
}

func (v *Viewer) serveGET(writer http.ResponseWriter, request *http.Request) {
	var formatter = v.formatterForRequest(request)
	var readCloser = v.filer.OpenReader(request)
	if err := v.filer.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)

		if filer.IsPathNotFoundError(err) {
			failer.ServeNotFound(writer, request)
			return
		}

		failer.ServeInternalServerError(writer, request)
		return
	}

	formatter.serveFromReader(readCloser, writer, request)
	readCloser.Close()
}

func (v *Viewer) formatterForRequest(request *http.Request) formatter {
	if mimeType := v.filer.MimeTypeForRequest(request); v.filer.Err() == nil {
		return mimeTypeFormatter(mimeType)
	}
	return newRawFormatter(fallbackMimeType)
}
