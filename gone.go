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
	"path"

	"github.com/fxnn/gone/authenticator"
	"github.com/fxnn/gone/config"
	"github.com/fxnn/gone/editor"
	"github.com/fxnn/gone/filestore"
	"github.com/fxnn/gone/router"
	"github.com/fxnn/gone/viewer"

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
	var htpasswdFilePath = htpasswdFilePath(contentRoot)

	var auth = authenticator.NewHttpBasicAuthenticator(htpasswdFilePath)
	var store = filestore.New(contentRoot, auth)

	var viewer = viewer.New(store)
	var editor = editor.New(store)
	var router = router.New(viewer, editor, auth)

	var handlerChain = context.ClearHandler(auth.AuthHandler(router))

	log.Fatal(http.ListenAndServe(cfg.BindAddress(), handlerChain))
}

func htpasswdFilePath(contentRootPath string) string {
	htpasswdFilePath := path.Join(contentRootPath, ".htpasswd")
	if _, err := os.Stat(htpasswdFilePath); err != nil && os.IsNotExist(err) {
		log.Printf("no .htpasswd found")
		return ""
	}
	log.Printf("using authentication data from .htpasswd")
	return htpasswdFilePath
}

func getwd() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	} else {
		panic(err)
	}
}
