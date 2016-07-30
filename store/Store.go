package store

import (
	"io"
	"net/http"
	"time"
)

// Store allows to read and write wiki contents and incorporates access
// control.
//
// Elementary to use the Store interface is the Err() method.
// As soon as an error occurs, all functions in Store turn to no-ops.
// The Err() method clears the error value and allows for error checking.
type Store interface {
	HasReadAccessForRequest(request *http.Request) bool
	HasWriteAccessForRequest(request *http.Request) bool
	HasDeleteAccessForRequest(request *http.Request) bool

	OpenReader(request *http.Request) io.ReadCloser
	OpenWriter(request *http.Request) io.WriteCloser

	ReadString(request *http.Request) string
	WriteString(request *http.Request, content string)

	Delete(request *http.Request)

	FileSizeForRequest(request *http.Request) int64
	MimeTypeForRequest(request *http.Request) string
	ModTimeForRequest(request *http.Request) time.Time

	// Err() clears and returns the error value.
	// It allows for error checking after one or more operations.
	Err() error
}
