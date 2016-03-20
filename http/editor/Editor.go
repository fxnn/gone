package editor

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fxnn/gone/http/failer"
	"github.com/fxnn/gone/http/router"
	"github.com/fxnn/gone/http/templates"
	"github.com/fxnn/gone/log"
	"github.com/fxnn/gone/store"
)

const (
	maxEditableBytes = 10 * 1024 * 1024 // 10 MiB of data
)

// The Editor is a HTTP Handler that serves the editor UI.
// While the UI itself is implemented in a HTML template, this type
// implements the logic behind the UI.
type Editor struct {
	store    store.Store
	renderer *templates.EditorRenderer
}

// New initializes a new instance ready to use.
// The instance includes a loaded and parsed template.
func New(l templates.Loader, s store.Store) *Editor {
	var renderer = templates.NewEditorRenderer()
	if err := renderer.LoadAndWatch(l); err != nil {
		panic(fmt.Errorf("couldn't load editor template: %s", err))
	}

	return &Editor{s, renderer}
}

func (e *Editor) isServeWriter(request *http.Request) bool {
	return request.Method == "POST"
}

func (e *Editor) isServeDeleter(request *http.Request) bool {
	return request.Method == "GET" && router.Is(router.ModeDelete, request)
}

func (e *Editor) isServeEditUI(request *http.Request) bool {
	return request.Method == "GET" && router.Is(router.ModeEdit, request)
}

func (e *Editor) isServeCreateUI(request *http.Request) bool {
	return request.Method == "GET" && router.Is(router.ModeCreate, request)
}

func (e *Editor) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if e.isServeWriter(request) {
		e.serveWriter(writer, request)
		return
	}

	if e.isServeDeleter(request) {
		e.serveDeleter(writer, request)
		return
	}

	if e.isServeCreateUI(request) || e.isServeEditUI(request) {
		e.serveEditUI(writer, request)
		return
	}

	log.Printf("%s %s: method not allowed", request.Method, request.URL)
	failer.ServeMethodNotAllowed(writer, request)
}

func (e *Editor) serveWriter(writer http.ResponseWriter, request *http.Request) {
	if !e.store.HasWriteAccessForRequest(request) {
		log.Printf("%s %s: no write permissions", request.Method, request.URL)
		failer.ServeUnauthorized(writer, request)
		return
	}

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

func (e *Editor) serveDeleter(writer http.ResponseWriter, request *http.Request) {
	if !e.store.HasDeleteAccessForRequest(request) {
		log.Printf("%s %s: no delete permissions", request.Method, request.URL)
		failer.ServeUnauthorized(writer, request)
		return
	}

	e.store.Delete(request)
	if err := e.store.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}
	log.Printf("%s %s: deleted", request.Method, request.URL)

	fmt.Fprintf(writer, "Successfully deleted")
}

func (e *Editor) serveEditUI(writer http.ResponseWriter, request *http.Request) {
	if !e.store.HasWriteAccessForRequest(request) {
		log.Printf("%s %s: no write permissions", request.Method, request.URL)
		failer.ServeUnauthorized(writer, request)
		return
	}
	if !router.Is(router.ModeCreate, request) && !e.store.HasReadAccessForRequest(request) {
		log.Printf("%s %s: no read permissions", request.Method, request.URL)
		failer.ServeUnauthorized(writer, request)
		return
	}

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
		} else if router.Is(router.ModeEdit, request) {
			log.Printf("%s %s: file to be edited does not exist: %s", request.Method, request.URL, err)
			failer.ServeNotFound(writer, request)
			return
		}
	} else if router.Is(router.ModeCreate, request) {
		log.Printf("%s %s: file to be created already exists: %s", request.Method, request.URL, err)
		failer.ServeConflict(writer, request)
		return
	}

	err := e.renderer.Render(writer, request.URL, content, router.Is(router.ModeEdit, request))
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	log.Printf("%s %s: served from template", request.Method, request.URL)
}

func (e *Editor) assertEditableTextFile(request *http.Request) error {
	bytes := e.store.FileSizeForRequest(request)
	if err := e.store.Err(); err != nil && store.IsPathNotFoundError(err) {
		// HINT: doesn't exist; that's pretty editable
		return nil
	}

	if bytes > maxEditableBytes {
		return fmt.Errorf(
			"the file size of %d bytes is larger than the allowed %d bytes",
			bytes, maxEditableBytes)
	}

	mimeType := e.store.MimeTypeForRequest(request)
	if e.store.Err() == nil && e.isKnownEditableMimeType(mimeType) {
		return nil
	}

	return fmt.Errorf("the mime type %s doesn't represent editable text", mimeType)
}

var knownEditableMimeTypePrefixes = []string{
	"application/xhtml+xml",
	"application/xml",
	"application/x-javascript",
	"application/x-latex",
	"application/x-sh",
	"application/x-tex",
	"application/x-troff",
	"text/",
}

func (e *Editor) isKnownEditableMimeType(mimeType string) bool {
	for _, prefix := range knownEditableMimeTypePrefixes {
		if strings.HasPrefix(mimeType, prefix) {
			return true
		}
	}
	return false
}
