package editor

import (
	"github.com/fxnn/gone/failer"
	"github.com/fxnn/gone/filer"
	"github.com/fxnn/gone/templates"
	"log"
	"net/http"
)

// The Editor is a HTTP Handler that serves the editor UI.
// While the UI itself is implemented in a HTML template, this type
// implements the logic behind the UI.
type Editor struct {
	filer    filer.Filer
	template templates.EditorTemplate
}

// Initializes a new instance ready to use.
// The instance includes a loaded and parsed template.
func New() Editor {
	var template = templates.LoadEditorTemplate()
	if err := template.Err(); err != nil {
		panic(err)
	}

	return Editor{filer.New(), template}
}

func (e *Editor) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		e.servePOST(writer, request)
		return
	}

	if request.Method == "GET" {
		e.serveGET(writer, request)
		return
	}

	log.Printf("%s %s: method not allowed", request.Method, request.URL)
	failer.ServeMethodNotAllowed(writer, request)
}

func (e *Editor) servePOST(writer http.ResponseWriter, request *http.Request) {
	var content = request.FormValue("content")
	if content == "" {
		log.Printf("%s %s: no valid content in request", request.Method, request.URL)
		failer.ServeBadRequest(writer, request)
		return
	}

	e.filer.WriteString(request, content)
	if err := e.filer.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}
	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, len(content))

	if request.FormValue("saveAndReturn") != "" {
		e.redirect(writer, request, request.URL.Path)
		return
	}

	e.redirect(writer, request, request.URL.Path+"?edit")
}

func (e *Editor) serveGET(writer http.ResponseWriter, request *http.Request) {
	var content = e.filer.ReadString(request)
	if err := e.filer.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		if filer.IsPathNotFoundError(err) {
			failer.ServeNotFound(writer, request)
		} else {
			failer.ServeInternalServerError(writer, request)
		}
		return
	}

	e.template.Render(writer, request.URL, content)
	if err := e.template.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	log.Printf("%s %s: served from template", request.Method, request.URL)
}

func (e *Editor) redirect(writer http.ResponseWriter, request *http.Request, location string) {
	http.Redirect(writer, request, location, http.StatusFound)
}
