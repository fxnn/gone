// Gone is a wiki written in Go. It is designed with server owners and
// administrators in mind and follows the KISS principles.
//
// Currently, gone simply starts a HTTP server at port 8080 and servers files
// from the working directory.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fxnn/gone/authenticator"
	"github.com/fxnn/gone/editor"
	"github.com/fxnn/gone/filer"
	"github.com/fxnn/gone/router"
	"github.com/fxnn/gone/viewer"

	"github.com/fxnn/gone/internal/github.com/gorilla/context"
)

func main() {
	var filer, authenticator = filerAndAuthenticator()

	var viewer = viewer.New(filer)
	var editor = editor.New(filer)
	var router = router.New(viewer, editor, authenticator)

	var handlerChain = context.ClearHandler(authenticator.AuthHandler(router))

	log.Fatal(http.ListenAndServe(":8080", handlerChain))
}

func filerAndAuthenticator() (f *filer.Filer, a *authenticator.HttpBasicAuthenticator) {
	f = filer.New(authenticator.NewNeverAuthenticated())
	f.SetContentRootPath(getwd())
	var htpasswdFilePath = f.HtpasswdFilePath()
	if err := f.Err(); err != nil {
		log.Printf("no .htpasswd found")
	} else {
		log.Printf("using authentication data from .htpasswd")
	}
	a = authenticator.NewHttpBasicAuthenticator(htpasswdFilePath)
	f.SetAuthenticator(a)
	return
}

func getwd() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	} else {
		panic(err)
	}
}
