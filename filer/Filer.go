package filer

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
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
	buf, f.err = ioutil.ReadAll(readCloser)
	readCloser.Close()
	return string(buf)
}

// Writes the given string into the given Writer and closes it.
func (f *Filer) writeAllAndClose(writeCloser io.WriteCloser, content string) {
	if f.err != nil {
		return
	}
	_, f.err = io.WriteString(writeCloser, content)
	writeCloser.Close()
}

// OpenReader opens a reader for the given request.
// A caller must close the reader after using it.
// Also, he must always check the Err() method.
func (f *Filer) OpenReader(request *http.Request) io.ReadCloser {
	return f.openReaderAtPath(f.pathFromRequest(request))
}

func (f *Filer) OpenWriter(request *http.Request) io.WriteCloser {
	return f.openWriterAtPath(f.pathFromRequest(request))
}

func (f *Filer) openReaderAtPath(p string) (reader io.ReadCloser) {
	if f.err != nil {
		return nil
	}
	reader, f.err = os.Open(p)
	f.wrapErr()
	return
}

func (f *Filer) openWriterAtPath(p string) (writer io.WriteCloser) {
	if f.err != nil {
		return nil
	}
	if !f.hasWriteAccessToPath(p) {
		f.err = NewAccessDeniedError(fmt.Sprintf("Access denied on %s", p))
		return nil
	}
	writer, f.err = os.Create(p)
	f.wrapErr()
	return
}

func (f *Filer) pathFromRequest(request *http.Request) string {
	var p = "." + request.URL.Path
	f.assertPathInsideWorkingDirectory(p)
	return p
}

func (f *Filer) HasWriteAccessForRequest(request *http.Request) bool {
	return f.hasWriteAccessToPath(f.pathFromRequest(request))
}

func (f *Filer) hasWriteAccessToPath(p string) bool {
	if f.err != nil {
		return false
	}
	info, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		// HINT: Inspect permissions of containing directory
		info, err = os.Stat(path.Dir(p))
	}
	if err != nil {
		f.err = err
		f.wrapErr()
		return false
	}
	return f.hasWriteAccessForFileMode(info.Mode())
}

func (f *Filer) hasWriteAccessForFileMode(mode os.FileMode) bool {
	// 0002 is the write permission write for others
	return mode&0002 != 0
}

// Wraps f.err to a filer-specific error, if possible
func (f *Filer) wrapErr() {
	if f.err != nil && os.IsNotExist(f.err) {
		if pathError, ok := f.err.(*os.PathError); ok {
			f.err = NewPathNotFoundError("path not found: " + pathError.Path)
		} else {
			f.err = NewPathNotFoundError(fmt.Sprintf("path not found: %s", f.err))
		}
	}
}
