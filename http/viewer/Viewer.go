package viewer

import (
	"net/http"
	"time"

	"github.com/fxnn/gone/http/failer"
	"github.com/fxnn/gone/http/templates"
	"github.com/fxnn/gone/log"
	"github.com/fxnn/gone/store"
)

// The Viewer serves HTTP requests with content from the filesystem.
type Viewer struct {
	store      store.Store
	formatters formatters
}

// New initializes a Viewer instance ready to use.
func New(l templates.Loader, s store.Store) *Viewer {
	return &Viewer{s, newFormatters(l)}
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
	if v.isNotModified(writer, request) {
		return
	}

	var formatter = v.formatterForRequest(request)
	var readCloser = v.store.OpenReader(request)
	if err := v.store.Err(); err != nil {
		v.serveError(writer, request, err)
		return
	}

	defer readCloser.Close()
	formatter.serveFromReader(readCloser, writer, request)
}

// isNotModified handles the complete Last-Modified / If-Modified-Since logic
// for HTTP caching.
func (v *Viewer) isNotModified(writer http.ResponseWriter, request *http.Request) bool {
	var modTime = v.store.ModTimeForRequest(request)

	if err := v.store.Err(); err == nil && !modTime.IsZero() {
		if ifModifiedSince, err := time.Parse(http.TimeFormat, request.Header.Get("If-Modified-Since")); err == nil {
			if modTime.Before(ifModifiedSince.Add(1*time.Second)) {
				writer.WriteHeader(http.StatusNotModified)
				return true
			}
		}

		writer.Header().Set("Last-Modified", modTime.UTC().Format(http.TimeFormat))
	}

	return false
}

func (v *Viewer) serveError(writer http.ResponseWriter, request *http.Request, err error) {
	v.log(request, err)

	if store.IsPathNotFoundError(err) {
		failer.ServeNotFound(writer, request)
		return
	}

	failer.ServeInternalServerError(writer, request)
}

func (v *Viewer) log(request *http.Request, err error) {
	log.Printf("%s %s: %s", request.Method, request.URL, err)
}

func (v *Viewer) formatterForRequest(request *http.Request) formatter {
	var mimeType = v.store.MimeTypeForRequest(request)
	v.store.Err() // don't care for errors
	return v.formatters.mimeTypeFormatter(mimeType)
}
