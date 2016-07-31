package store

import (
	"mime"

	"github.com/fxnn/gone/log"
)

const (
	FallbackMimeType = "application/octet-stream"
	MarkdownMimeType = "text/markdown"
	UrlMimeType = "text/url"
)

func init() {
	// http://superuser.com/a/285878
	registerMarkdownExtension("md")
	registerMarkdownExtension("mkdn")
	registerMarkdownExtension("markdown")
	registerMarkdownExtension("mdown")
	registerMarkdownExtension("mkd")
	registerMarkdownExtension("mdwn")
	registerMarkdownExtension("mdtxt")
	registerMarkdownExtension("mdtext")
	registerMarkdownExtension("text")

	registerExtension("url", UrlMimeType)
}

func registerMarkdownExtension(ext string) {
	registerExtension(ext, MarkdownMimeType)
}

func registerExtension(ext string, mimeType string) {
	if err := mime.AddExtensionType("."+ext, mimeType); err != nil {
		log.Printf("error adding %s as extension for Mime type %s: %s", ext, mimeType, err)
	}
}
