package filer

import (
	"io"
	"io/ioutil"
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

// Returns the requested content as string.
// A caller must always check the Err() method.
func (f *Filer) ReadString(request *http.Request) string {
	if f.err != nil {
		return ""
	}
	return f.readAllAndClose(f.OpenReader(request))
}

// Reads everything into the given Reader until EOF and closes it.
func (f *Filer) readAllAndClose(readCloser io.ReadCloser) (result string) {
	if f.err != nil {
		return ""
	}
	var buf []byte
	buf, f.err = ioutil.ReadAll(readCloser)
	readCloser.Close()
	return string(buf)
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
