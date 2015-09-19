// Gone is a wiki written in Go. It is designed with server owners and
// administrators in mind and follows the KISS principles.
//
// Currently, gone simply starts a HTTP server at port 8080 and servers files
// from the working directory.
package main

import (
	"github.com/fxnn/gone/editor"
	"github.com/fxnn/gone/filer"
	"log"
	"net/http"
)

func main() {
	http.Handle("/view/", http.StripPrefix("/view", filer.NewHandler()))
	http.Handle("/edit/", must(editor.NewHandler()))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func must(handler http.Handler, err error) http.Handler {
	if err != nil {
		log.Panic(err)
	}
	return handler
}
