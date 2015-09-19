// Gone is a wiki written in Go. It is designed with server owners and
// administrators in mind and follows the KISS principles.
//
// Currently, gone simply starts a HTTP server at port 8080 and servers files
// from the working directory.
package main

import (
	"github.com/fxnn/gone/filer"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", filer.NewHandler())
}
