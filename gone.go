// Gone is a wiki written in Go. It is designed with server owners and
// administrators in mind and follows the KISS principles.
//
// Currently, gone simply starts a HTTP server at port 8080 and servers files
// from the working directory.
package main

import (
	"github.com/fxnn/gone/editor"
	"github.com/fxnn/gone/router"
	"github.com/fxnn/gone/viewer"
	"log"
	"net/http"
)

func main() {
	var viewer = viewer.New()
	var editor = editor.New()
	var router = router.New(&viewer, &editor)

	log.Fatal(http.ListenAndServe(":8080", &router))
}
