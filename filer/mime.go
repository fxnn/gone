package filer

import (
	"log"
	"mime"
)

const (
	MarkdownMimeType = "text/markdown"
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
}

func registerMarkdownExtension(ext string) {
	if err := mime.AddExtensionType("."+ext, MarkdownMimeType); err != nil {
		log.Printf("error adding %s as markdown extension: %s", ext, err)
	}
}
