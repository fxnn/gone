package viewer

import (
	"github.com/fxnn/gone/failer"
	"github.com/fxnn/gone/store"
	"log"
	"net/http"
)

// The Viewer serves HTTP requests with content from the filesystem.
type Viewer struct {
	store store.Store
}

// New initializes a Viewer instance ready to use.
func New(s store.Store) *Viewer {
	return &Viewer{s}
}

func (v *Viewer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !v.store.HasReadAccessForRequest(request) {
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
	var readCloser = v.store.OpenReader(request)
	if err := v.store.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)

		if store.IsPathNotFoundError(err) {
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
	var mimeType = v.store.MimeTypeForRequest(request)
	v.store.Err() // don't care for errors
	return mimeTypeFormatter(mimeType)
}
