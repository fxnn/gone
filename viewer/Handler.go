package viewer

import (
	"github.com/fxnn/gone/failer"
	"github.com/fxnn/gone/filer"
	"io"
	"log"
	"net/http"
)

// The `Handler` in this package serves HTTP requests with content from the
// filesystem.
type Handler struct {
	filer filer.Filer
}

// Initializes a zeroe'd instance ready to use.
func New() Handler {
	return Handler{filer.New()}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		h.serveNonGET(writer, request)
		return
	}

	h.serveGET(writer, request)
}

func (h *Handler) serveNonGET(writer http.ResponseWriter, request *http.Request) {
	failer.ServeMethodNotAllowed(writer, request)
}

func (h *Handler) serveGET(writer http.ResponseWriter, request *http.Request) {
	var readCloser = h.filer.OpenReader(request)
	if err := h.filer.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		if filer.IsPathNotFoundError(err) {
			failer.ServeNotFound(writer, request)
		} else {
			failer.ServeInternalServerError(writer, request)
		}
		return
	}

	h.serveFromReader(readCloser, writer, request)
	readCloser.Close()
}

func (h *Handler) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	written, err := io.Copy(writer, reader)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, written)
}
