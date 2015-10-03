package router

import (
	"github.com/fxnn/gone/failer"
	"log"
	"net/http"
)

type Router struct {
	editor        http.Handler
	viewer        http.Handler
	authenticator http.Handler
}

func New(viewer http.Handler, editor http.Handler, authenticator http.Handler) Router {
	return Router{editor, viewer, authenticator}
}

func (r Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err = request.ParseForm()
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeBadRequest(writer, request)
	} else if IsModeLogin(request) {
		r.authenticator.ServeHTTP(writer, request)
	} else if IsModeEdit(request) || IsModeCreate(request) {
		r.editor.ServeHTTP(writer, request)
	} else {
		r.viewer.ServeHTTP(writer, request)
	}
}
