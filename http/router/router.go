package router

import (
	"net/http"

	"github.com/fxnn/gone/http/failer"
	"github.com/fxnn/gone/log"
)

// Router encapsulates http.Handler instances for all relevant views and
// invokes the right one for each request.
type Router struct {
	editor        http.Handler
	viewer        http.Handler
	authenticator http.Handler
}

// New constructs a new instance ready to use.
func New(viewer http.Handler, editor http.Handler, authenticator http.Handler) *Router {
	return &Router{editor, viewer, authenticator}
}

func (r Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err = request.ParseForm()
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err)
		failer.ServeBadRequest(writer, request)
	} else if Is(ModeLogin, request) {
		r.authenticator.ServeHTTP(writer, request)
	} else if Is(ModeEdit, request) || Is(ModeCreate, request) || Is(ModeDelete, request) {
		r.editor.ServeHTTP(writer, request)
	} else {
		r.viewer.ServeHTTP(writer, request)
	}
}
