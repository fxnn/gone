// Gone is a wiki written in Go. It is designed with server owners and
// administrators in mind and follows the KISS principles.
//
// Currently, gone simply starts a HTTP server at port 8080 and servers files
// from the working directory.
package main

import (
	"log"
	"net/http"

	"github.com/fxnn/gone/authenticator"
	"github.com/fxnn/gone/editor"
	"github.com/fxnn/gone/filer"
	"github.com/fxnn/gone/router"
	"github.com/fxnn/gone/viewer"

	"github.com/gorilla/context"
)

func main() {
	var authenticator = authenticator.NewHttpBasicAuthenticator()
	var filer = filer.New(authenticator)

	var viewer = viewer.New(filer)
	var editor = editor.New(filer)
	var router = router.New(viewer, editor, authenticator)

	var handlerChain = context.ClearHandler(authenticator.AuthHandler(router))

	log.Fatal(http.ListenAndServe(":8080", handlerChain))
}
