package editor

import (
	"github.com/fxnn/gone/filer"
	"github.com/fxnn/gone/templates"
	"io"
	"log"
	"net/http"
)

// Serves the editor UI.
type Handler struct {
	filer    filer.Filer
	template templates.EditorTemplate
}

// Initializes a zeroe'd instance ready to use.
func New() Handler {
	var template = templates.LoadEditorTemplate()
	if err := template.Err(); err != nil {
		panic(err)
	}

	return Handler{filer.New(), template}
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
	var content = h.filer.ReadString(request)
	if err := h.filer.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		if filer.IsPathNotFoundError(err) {
			h.serveNotFound(writer, request)
		} else {
			h.serveInternalServerError(writer, request)
		}
		return
	}

	h.template.Render(writer, request.URL, content)
	if err := h.template.Err(); err != nil {
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
