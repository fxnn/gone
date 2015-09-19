package editor

import (
	"github.com/fxnn/gone/filer"
	"io"
	"log"
	"net/http"
)

// Serves the editor UI.
type Handler struct {
	filer        filer.Filer
	editTemplate editTemplate
}

// Initializes a zeroe'd instance ready to use.
func New() (*Handler, error) {
	var editTemplate, err = loadEditTemplate()
	if err != nil {
		return nil, err
	}

	return &Handler{filer.New(), editTemplate}, nil
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		h.serveNonGET(writer, request)
		return
	}

	h.serveGET(writer, request)
}

func (h *Handler) serveNonGET(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusMethodNotAllowed)
	io.WriteString(writer, "Oops, method not allowed")
}

func (h *Handler) serveGET(writer http.ResponseWriter, request *http.Request) {
	var data = make(map[string]string)
	data["path"] = request.URL.Path
	data["content"] = h.filer.ReadString(request)
	if h.filer.Err() != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, h.filer.Err())
		if filer.IsPathNotFoundError(h.filer.Err()) {
			h.serveNotFound(writer, request)
		} else {
			h.serveInternalServerError(writer, request)
		}
		return
	}

	var err = h.editTemplate.Execute(writer, data)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		h.serveInternalServerError(writer, request)
		return
	}

	log.Printf("%s %s: served from template", request.Method, request.URL)
}

func (h *Handler) serveNotFound(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	io.WriteString(writer, "Oops, file not found")
}

func (h *Handler) serveInternalServerError(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)
	io.WriteString(writer, "Oops, internal server error")
}
