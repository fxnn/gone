package viewer

import (
	"github.com/fxnn/gone/failer"
	"github.com/fxnn/gone/filer"
	"io"
	"log"
	"net/http"
)

// The Viewer serves HTTP requests with content from the filesystem.
type Viewer struct {
	filer filer.Filer
}

// Initializes a zeroe'd instance ready to use.
func New() Viewer {
	return Viewer{filer.New()}
}

func (v *Viewer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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
	var readCloser = v.filer.OpenReader(request)
	if err := v.filer.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		if filer.IsPathNotFoundError(err) {
			failer.ServeNotFound(writer, request)
		} else {
			failer.ServeInternalServerError(writer, request)
		}
		return
	}

	v.serveFromReader(readCloser, writer, request)
	readCloser.Close()
}

func (v *Viewer) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	written, err := io.Copy(writer, reader)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, written)
}
