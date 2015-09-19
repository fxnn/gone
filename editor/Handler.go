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
func NewHandler() (*Handler, error) {
	var editTemplate, err = loadEditTemplate()
	if err != nil {
		return nil, err
	}

	return &Handler{filer.NewFiler(), editTemplate}, nil
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
	var err = h.editTemplate.Execute(writer, "no data yet")
	if err != nil {
		h.serveInternalServerError(writer, request)
		log.Printf("%s %s: %s", request.Method, request.URL, err.Error())
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
