package filer

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fxnn/gone/authenticator"
)

// Maps incoming HTTP requests to the file system.
type Filer struct {
	accessControl
}

// Initializes a zeroe'd instance ready to use.
func New(authenticator authenticator.Authenticator) *Filer {
	return &Filer{newAccessControl(authenticator)}
}

// Returns the requested content as string.
// A caller must always check the Err() method.
func (f *Filer) ReadString(request *http.Request) string {
	if f.err != nil {
		return ""
	}
	return f.readAllAndClose(f.OpenReader(request))
}

// Writes the given content into a file pointed to by the request.
// A caller must always check the Err() method.
func (f *Filer) WriteString(request *http.Request, content string) {
	if f.err != nil {
		return
	}
	f.writeAllAndClose(f.OpenWriter(request), content)
}

// Reads everything into the given Reader until EOF and closes it.
func (f *Filer) readAllAndClose(readCloser io.ReadCloser) (result string) {
	if f.err != nil {
		return ""
	}
	var buf []byte
	buf, err := ioutil.ReadAll(readCloser)
	f.setErr(err)
	readCloser.Close()
	return string(buf)
}

// Writes the given string into the given Writer and closes it.
func (f *Filer) writeAllAndClose(writeCloser io.WriteCloser, content string) {
	if f.err != nil {
		return
	}
	_, err := io.WriteString(writeCloser, content)
	f.setErr(err)
	writeCloser.Close()
}

// OpenReader opens a reader for the given request.
// A caller must close the reader after using it.
// Also, he must always check the Err() method.
func (f *Filer) OpenReader(request *http.Request) io.ReadCloser {
	f.assertHasReadAccessForRequest(request)
	if f.err != nil {
		return nil
	}
	return f.openReaderAtPath(f.pathFromRequest(request))
}

func (f *Filer) OpenWriter(request *http.Request) io.WriteCloser {
	f.assertHasWriteAccessForRequest(request)
	if f.err != nil {
		return nil
	}
	return f.openWriterAtPath(f.pathFromRequest(request))
}

func (f *Filer) openReaderAtPath(p string) (reader io.ReadCloser) {
	if f.err != nil {
		return nil
	}
	reader, err := os.Open(p)
	f.setErr(err)
	return
}

func (f *Filer) openWriterAtPath(p string) (writer io.WriteCloser) {
	if f.err != nil {
		return nil
	}
	writer, err := os.Create(p)
	f.setErr(err)
	return
}
