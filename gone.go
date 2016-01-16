// Gone is a wiki written in Go. It is designed with server owners and
// administrators in mind and follows the KISS principles.
//
// By default, gone simply starts a HTTP server at port 8080 and servers files
// from the working directory.
// Invoke with -help flag to see configuration options.
package main

import (
	"log"
	"os"

	"github.com/fxnn/gone/authenticator"
	"github.com/fxnn/gone/config"
	"github.com/fxnn/gone/http"
	"github.com/fxnn/gone/http/templates"
	"github.com/fxnn/gone/store/filestore"
	"github.com/fxnn/gopath"
)

const defaultTemplateDirectoryName = ".templates"

func main() {
	cfg := config.FromCommandline()

	switch cfg.Command {
	case config.CommandExportTemplates:
		exportTemplates(cfg)
	case config.CommandListen:
		listen(cfg)
	case config.CommandHelp:
		config.PrintUsage()
	}
}

func exportTemplates(cfg config.Config) {
	var target = templatePath(contentRoot(), cfg)
	if target.IsEmpty() {
		target = contentRoot().JoinPath(defaultTemplateDirectoryName)
	}

	if err := templates.NewStaticLoader().WriteAllTemplates(target); err != nil {
		log.Fatalf("error exporting templates: %s", err)
	}
}

func listen(cfg config.Config) {
	var cr = contentRoot()

	var auth = createAuthenticator(cr, cfg)
	var store = filestore.New(cr, auth)
	var loader = createLoader(cr, cfg)

	http.ListenAndServe(cfg.BindAddress, auth, store, loader)
}

func createAuthenticator(
	contentRoot gopath.GoPath,
	cfg config.Config,
) *authenticator.HttpBasicAuthenticator {
	var htpasswdFile = htpasswdFilePath(contentRoot)
	if cfg.RequireSSLHeader != "" {
		log.Printf("Requiring SSL header %s on login", cfg.RequireSSLHeader)
	}
	return authenticator.NewHttpBasicAuthenticator(
		htpasswdFile,
		cfg.RequireSSLHeader,
		cfg.BruteforceMaxDelay,
		cfg.BruteforceDelayStep,
		cfg.BruteforceDelayStep/5,
		cfg.BruteforceDelayStep/20,
		cfg.BruteforceDropDelayAfter,
	)
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
	var templatePath = templatePath(contentRoot, cfg)

	if !templatePath.IsEmpty() {
		return templates.NewFilesystemLoader(templatePath)
	}

	return templates.NewStaticLoader()
}

func templatePath(contentRoot gopath.GoPath, cfg config.Config) (result gopath.GoPath) {
	// configuration
	result = gopath.FromPath(cfg.TemplatePath)
	if !result.IsEmpty() {
		if !result.IsDirectory() {
			log.Fatalf("configured template path is no directory: %s", result.Path())
		}
		log.Printf("using templates from %s (by configuration)", result.Path())
		return result
	}

	// convention
	result = contentRoot.JoinPath(defaultTemplateDirectoryName)
	if result.IsDirectory() {
		log.Printf("using templates from %s (by convention)", result.Path())
		return result
	}

	// default
	log.Printf("using default templates")
	return gopath.Empty()
}

func contentRoot() gopath.GoPath {
	if wd, err := os.Getwd(); err == nil {
		return gopath.FromPath(wd)
	} else {
		panic(err)
	}
}
