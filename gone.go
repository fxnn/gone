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
	"github.com/fxnn/gone/config"
	"github.com/fxnn/gone/editor"
	"github.com/fxnn/gone/filestore"
	"github.com/fxnn/gone/router"
	"github.com/fxnn/gone/templates"
	"github.com/fxnn/gone/viewer"
	"github.com/fxnn/gopath"

	"github.com/gorilla/context"
)

func main() {
	cfg := config.FromCommandline()

	switch cfg.Command() {
	case config.CommandListen:
		listen(cfg)
	case config.CommandHelp:
		config.PrintUsage()
	}
}

func listen(cfg config.Config) {
	var contentRoot = getwd()

	var auth = createAuthenticator(contentRoot)
	var store = filestore.New(contentRoot, auth)
	var loader = templates.NewStaticLoader()

	var viewer = viewer.New(loader, store)
	var editor = editor.New(loader, store)
	var router = router.New(viewer, editor, auth)

	var handlerChain = context.ClearHandler(auth.AuthHandler(router))

	log.Fatal(http.ListenAndServe(cfg.BindAddress(), handlerChain))
}

func createAuthenticator(contentRoot gopath.GoPath) *authenticator.HttpBasicAuthenticator {
	var htpasswdFile = htpasswdFilePath(contentRoot)
	return authenticator.NewHttpBasicAuthenticator(htpasswdFile)
}

func htpasswdFilePath(contentRoot gopath.GoPath) gopath.GoPath {
	htpasswdFile := contentRoot.JoinPath(".htpasswd")
	if !htpasswdFile.IsExists() {
		log.Printf("no .htpasswd found")
	} else {
		log.Printf("using authentication data from .htpasswd")
	}
	return htpasswdFile
}

func getwd() gopath.GoPath {
	if wd, err := os.Getwd(); err == nil {
		return gopath.FromPath(wd)
	} else {
		panic(err)
	}
}
