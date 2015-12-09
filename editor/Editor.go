package editor

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fxnn/gone/failer"
	"github.com/fxnn/gone/router"
	"github.com/fxnn/gone/store"
	"github.com/fxnn/gone/templates"
)

const (
	maxEditableBytes = 10 * 1024 * 1024 // 10 MiB of data
)

// The Editor is a HTTP Handler that serves the editor UI.
// While the UI itself is implemented in a HTML template, this type
// implements the logic behind the UI.
type Editor struct {
	store    store.Store
	template templates.EditorTemplate
}

// Initializes a new instance ready to use.
// The instance includes a loaded and parsed template.
func New(s store.Store) *Editor {
	var template = templates.LoadEditorTemplate()
	if err := template.Err(); err != nil {
		panic(err)
	}

	return &Editor{s, template}
}

func (e *Editor) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if !e.store.HasWriteAccessForRequest(request) {
		log.Printf("%s %s: no write permissions", request.Method, request.URL)
		failer.ServeUnauthorized(writer, request)
		return
	}

	if request.Method == "POST" {
		e.servePOST(writer, request)
		return
	}

	if !e.store.HasReadAccessForRequest(request) {
		log.Printf("%s %s: no read permissions", request.Method, request.URL)
		failer.ServeUnauthorized(writer, request)
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

	e.store.WriteString(request, content)
	if err := e.store.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}
	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, len(content))

	if request.FormValue("saveAndReturn") != "" {
		router.RedirectToViewMode(writer, request)
		return
	}

	router.RedirectToEditMode(writer, request)
}

func (e *Editor) serveGET(writer http.ResponseWriter, request *http.Request) {
	if err := e.assertEditableTextFile(request); err != nil {
		log.Printf("%s %s: no editable text file: %s", request.Method, request.URL, err)
		failer.ServeUnsupportedMediaType(writer, request)
		return
	}

	var content = e.store.ReadString(request)
	if err := e.store.Err(); err != nil {
		if !store.IsPathNotFoundError(err) {
			log.Printf("%s %s: %s", request.Method, request.URL, err)
			failer.ServeInternalServerError(writer, request)
			return
		} else if router.IsModeEdit(request) {
			log.Printf("%s %s: file to be edited does not exist: %s", request.Method, request.URL, err)
			failer.ServeNotFound(writer, request)
			return
		}
	} else if router.IsModeCreate(request) {
		log.Printf("%s %s: file to be created already exists: %s", request.Method, request.URL, err)
		failer.ServeConflict(writer, request)
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

func (e *Editor) assertEditableTextFile(request *http.Request) error {
	bytes := e.store.FileSizeForRequest(request)
	if err := e.store.Err(); err != nil && os.IsNotExist(err) {
		// HINT: doesn't exist; that's pretty editable
		return nil
	}

	if bytes > maxEditableBytes {
		return fmt.Errorf(
			"the file size of %d bytes is larger than the allowed %d bytes",
			bytes, maxEditableBytes)
	}

	mimeType := e.store.MimeTypeForRequest(request)
	if e.store.Err() == nil && strings.HasPrefix(mimeType, "text/") {
		return nil
	}

	return fmt.Errorf("the mime type %s doesn't represent editable text", mimeType)
}
