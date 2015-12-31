// Gone is a wiki written in Go. It is designed with server owners and
// administrators in mind and follows the KISS principles.
//
// By default, gone simply starts a HTTP server at port 8080 and servers files
// from the working directory.
// Invoke with -help flag to see configuration options.
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

const defaultTemplateDirectoryName = ".templates"

func main() {
	cfg := config.FromCommandline()

	switch cfg.Command() {
	case config.CommandExportTemplates:
		exportTemplates(cfg)
	case config.CommandListen:
		listen(cfg)
	case config.CommandHelp:
		config.PrintUsage()
	}
}

func exportTemplates(cfg config.Config) {
	var target = getwd().JoinPath(defaultTemplateDirectoryName)
	if err := templates.NewStaticLoader().WriteAllTemplates(target); err != nil {
		log.Fatalf("error exporting templates: %s", err)
	}
}

func listen(cfg config.Config) {
	var contentRoot = getwd()

	var auth = createAuthenticator(contentRoot)
	var store = filestore.New(contentRoot, auth)
	var loader = createLoader(contentRoot, cfg)

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

func createLoader(contentRoot gopath.GoPath, cfg config.Config) templates.Loader {
	var templatePath gopath.GoPath

	// configuration
	templatePath = gopath.FromPath(cfg.TemplatePath())
	if !templatePath.IsEmpty() {
		if !templatePath.IsDirectory() {
			log.Fatalf("configured template path is no directory: %s", templatePath.Path())
		}
		log.Printf("using templates from %s (by configuration)", templatePath.Path())
		return templates.NewFilesystemLoader(templatePath)
	}

	// convention
	templatePath = contentRoot.JoinPath(defaultTemplateDirectoryName)
	if templatePath.IsDirectory() {
		log.Printf("using templates from %s (by convention)", templatePath.Path())
		return templates.NewFilesystemLoader(templatePath)
	}

	// default
	log.Printf("using default templates")
	return templates.NewStaticLoader()
}

func getwd() gopath.GoPath {
	if wd, err := os.Getwd(); err == nil {
		return gopath.FromPath(wd)
	} else {
		panic(err)
	}
}
