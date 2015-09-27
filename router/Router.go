package router

import (
	"github.com/fxnn/gone/failer"
	"log"
	"net/http"
)

type Router struct {
	editor http.Handler
	viewer http.Handler
}

func New(viewer http.Handler, editor http.Handler) Router {
	return Router{editor, viewer}
}

func (r Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err = request.ParseForm()
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeBadRequest(writer, request)
		return
	}
	if _, ok := request.Form["edit"]; ok {
		r.editor.ServeHTTP(writer, request)
		return
	}
	r.viewer.ServeHTTP(writer, request)
}
