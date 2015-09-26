package router

import (
	"github.com/fxnn/gone/failer"
	"log"
	"net/http"
)

type Handler struct {
	editor http.Handler
	viewer http.Handler
}

func New(viewer http.Handler, editor http.Handler) Handler {
	return Handler{editor, viewer}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err = request.ParseForm()
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeBadRequest(writer, request)
		return
	}
	if _, ok := request.Form["edit"]; ok {
		h.editor.ServeHTTP(writer, request)
		return
	}
	h.viewer.ServeHTTP(writer, request)
}
