package filer

import (
	"io"
	"net/http"
	"os"
)

// Maps incoming HTTP requests to the file system.
type Filer struct {
	basicFiler
}

// Initializes a zeroe'd instance ready to use.
func New() Filer {
	return Filer{}
}

// OpenReader opens a reader for the given request.
// A caller must close the reader after using it.
// Also, he must always check the Err() method.
func (f *Filer) OpenReader(request *http.Request) io.ReadCloser {
	var path = "." + request.URL.Path

	f.assertPathInsideWorkingDirectory(path)
	return f.openReaderAtPath(path)
}

func (f *Filer) openReaderAtPath(path string) (reader io.ReadCloser) {
	if f.err != nil {
		return nil
	}
	reader, f.err = os.Open(path)
	return
}
