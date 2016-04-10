package uploader

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fxnn/gone/http/failer"
	"github.com/fxnn/gone/http/router"
	"github.com/fxnn/gone/http/templates"
	"github.com/fxnn/gone/store"
)

// The Uploader is a HTTP handler that servers the uploader UI.
type Uploader struct {
	store    store.Store
	renderer *templates.UploaderRenderer
}

// New initializes a new instance ready to use.
// The instance includes a loaded and parsed template.
func New(l templates.Loader, s store.Store) *Uploader {
	var renderer = templates.NewUploaderRenderer()
	if err := renderer.LoadAndWatch(l); err != nil {
		panic(fmt.Errorf("couldn't load uploader template: %s", err))
	}

	return &Uploader{s, renderer}
}

func (u *Uploader) isServeUpload(request *http.Request) bool {
	return request.Method == "POST"
}

func (u *Uploader) isServeUI(request *http.Request) bool {
	return request.Method == "GET"
}

func (u *Uploader) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if u.isServeUpload(request) {
		u.serveUpload(writer, request)
		return
	}

	if u.isServeUI(request) {
		u.serveUI(writer, request)
		return
	}

	log.Printf("%s %s: method not allowed", request.Method, request.URL)
	failer.ServeMethodNotAllowed(writer, request)
}

func (u *Uploader) serveUpload(writer http.ResponseWriter, request *http.Request) {
	if !u.store.HasWriteAccessForRequest(request) {
		log.Printf("%s %s: no write permissions", request.Method, request.URL)
		failer.ServeUnauthorized(writer, request)
		return
	}

	var content = "TODO"
	if content == "" {
		log.Printf("%s %s: no valid content in request", request.Method, request.URL)
		failer.ServeBadRequest(writer, request)
		return
	}

	u.store.WriteString(request, content)
	if err := u.store.Err(); err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}
	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, len(content))

	router.RedirectToViewMode(writer, request)
}

func (u *Uploader) serveUI(writer http.ResponseWriter, request *http.Request) {
	if !u.store.HasWriteAccessForRequest(request) {
		log.Printf("%s %s: no write permissions", request.Method, request.URL)
		failer.ServeUnauthorized(writer, request)
		return
	}

	targetAlreadyExists := u.store.FileForRequestExists(request)
	err := u.renderer.Render(writer, request.URL, targetAlreadyExists)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}

	log.Printf("%s %s: served from template", request.Method, request.URL)
}
