package editor

import (
	"github.com/fxnn/gone/failer"
	"github.com/fxnn/gone/filer"
	"github.com/fxnn/gone/templates"
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
	if request.Method == "POST" {
		h.servePOST(writer, request)
		return
	}

	if request.Method == "GET" {
		h.serveGET(writer, request)
		return
	}

	log.Printf("%s %s: method not allowed", request.Method, request.URL)
	failer.ServeMethodNotAllowed(writer, request)
}

func (h *Handler) servePOST(writer http.ResponseWriter, request *http.Request) {
	var content = request.FormValue("content")
	if content == "" {
		log.Printf("%s %s: no valid content in request", request.Method, request.URL)
		failer.ServeBadRequest(writer, request)
		return
	}

	h.filer.WriteString(request, content)
	if err := h.filer.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}
	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, len(content))

	if request.FormValue("saveAndReturn") != "" {
		h.redirect(writer, request, request.URL.Path)
		return
	}

	h.redirect(writer, request, request.URL.Path+"?edit")
}

func (h *Handler) serveGET(writer http.ResponseWriter, request *http.Request) {
	var content = h.filer.ReadString(request)
	if err := h.filer.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		if filer.IsPathNotFoundError(err) {
			failer.ServeNotFound(writer, request)
		} else {
			failer.ServeInternalServerError(writer, request)
		}
		return
	}

	h.template.Render(writer, request.URL, content)
	if err := h.template.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	log.Printf("%s %s: served from template", request.Method, request.URL)
}

func (h *Handler) redirect(writer http.ResponseWriter, request *http.Request, location string) {
	http.Redirect(writer, request, location, http.StatusFound)
}
