package templates

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/fxnn/gone/http/failer"
	"github.com/fxnn/gone/log"
)

type TemplateDeliverer struct {
	loader Loader
}

func NewTemplateDeliverer(l Loader) *TemplateDeliverer {
	return &TemplateDeliverer{l}
}

func (e *TemplateDeliverer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		log.Printf("%s %s: wrong method for template handling", request.Method, request.URL)
		failer.ServeMethodNotAllowed(writer, request)
		return
	}

	readCloser, err := e.loader.LoadResource(request.URL.Path)
	if err != nil {
		log.Printf("%s %s: error while opening template resource: %s",
			request.Method, request.URL, err)
		failer.ServeNotFound(writer, request)
		return
	}
	defer readCloser.Close()

	writer.Header().Set("Content-Type", e.detectMimeType(request))

	_, err = io.Copy(writer, readCloser)
	if err != nil {
		log.Printf("%s %s: error while writing template resource to output: %s",
			request.Method, request.URL, err)
		failer.ServeInternalServerError(writer, request)
		return
	}
}

func (e *TemplateDeliverer) detectMimeType(r *http.Request) string {
	path := r.URL.Path
	ext := filepath.Ext(path)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "application/octet-stream"
	}
	return mimeType
}
