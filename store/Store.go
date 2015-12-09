package store

import (
	"io"
	"net/http"
)

type Store interface {
	HasReadAccessForRequest(request *http.Request) bool
	HasWriteAccessForRequest(request *http.Request) bool

	OpenReader(request *http.Request) io.ReadCloser
	OpenWriter(request *http.Request) io.WriteCloser

	ReadString(request *http.Request) string
	WriteString(request *http.Request, content string)

	FileSizeForRequest(request *http.Request) int64
	MimeTypeForRequest(request *http.Request) string

	Err() error
}
