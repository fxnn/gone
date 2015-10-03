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
	"github.com/fxnn/gone/router"
	"github.com/fxnn/gone/viewer"
)

func main() {
	var viewer = viewer.New()
	var editor = editor.New()
	var authenticator = authenticator.New()
	var router = router.New(&viewer, &editor, &authenticator)

	log.Fatal(http.ListenAndServe(":8080", &router))
}
