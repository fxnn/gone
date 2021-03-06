package http

import (
	"net/http"

	"github.com/fxnn/gone/authenticator"
	"github.com/fxnn/gone/http/editor"
	"github.com/fxnn/gone/http/router"
	"github.com/fxnn/gone/http/templates"
	"github.com/fxnn/gone/http/viewer"
	"github.com/fxnn/gone/log"
	"github.com/fxnn/gone/store"

	"github.com/gorilla/context"
)

// ListenAndServe brings up the web server component, waits for incoming HTTP
// requests on the given bindAddress and serves them.
func ListenAndServe(
	bindAddress string,
	auth authenticator.HttpAuthenticator,
	store store.Store,
	loader templates.Loader) {
	var templateDeliverer = templates.NewTemplateDeliverer(loader)
	var viewer = viewer.New(loader, store)
	var editor = editor.New(loader, store)
	var router = router.New(viewer, editor, templateDeliverer, auth.LoginHandler())

	var handlerChain = RequestLogger(
		context.ClearHandler(
			auth.MiddlewareHandler(
				router)))

	log.Fatal(http.ListenAndServe(bindAddress, handlerChain))
}
