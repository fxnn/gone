package filestore

import (
	"io"
	"net/http"
	"os"

	"github.com/fxnn/gone/authenticator"
	"github.com/fxnn/gone/store"
)

// fileStore implements the storage using plain file system.
type fileStore struct {
	*errStore
	*ioUtil
	*pathIO
	*mimeDetector
	*accessControl
}

// New initializes a zeroe'd instance ready to use.
func New(contentRoot string, authenticator authenticator.Authenticator) store.Store {
	var s = newErrStore()
	var i = newIOUtil(s)
	var p = newPathIO(contentRoot, s)
	var m = newMimeDetector(p, s)
	var a = newAccessControl(authenticator, p, s)
	return &fileStore{s, i, p, m, a}
}

// Err returns and clears the recorder error.
//
// As soon as an error inside the filestore occurs, all operations turn into
// no-ops.
// Use this method to regularly check for errors.
func (f *fileStore) Err() error {
	return f.errAndClear()
}

func (f *fileStore) MimeTypeForRequest(request *http.Request) string {
	if f.hasErr() {
		return ""
	}
	return f.mimeTypeForPath(f.pathFromRequest(request))
}

// FileSizeForRequest returns the size of the underlying file in bytes, if any,
// or sets the Err() value.
func (f *fileStore) FileSizeForRequest(request *http.Request) int64 {
	p := f.pathFromRequest(request).Stat()
	if p.HasErr() {
		f.setErr(p.Err())
		return -1
	}

	return p.FileInfo().Size()
}

// ReadString returns the requested content as string.
// A caller must always check the Err() method.
func (f *fileStore) ReadString(request *http.Request) string {
	if f.hasErr() {
		return ""
	}
	return f.readAllAndClose(f.OpenReader(request))
}

// WriteString writes the given content into a file pointed to by the request.
// A caller must always check the Err() method.
func (f *fileStore) WriteString(request *http.Request, content string) {
	if f.hasErr() {
		return
	}
	f.writeAllAndClose(f.OpenWriter(request), content)
}

// OpenReader opens a reader for the given request.
// A caller must close the reader after using it.
// Also, he must always check the Err() method.
//
// The method handles access control.
func (f *fileStore) OpenReader(request *http.Request) io.ReadCloser {
	if f.hasErr() {
		return nil
	}
	f.assertHasReadAccessForRequest(request)
	return f.openReaderAtPath(f.pathFromRequest(request))
}

// OpenWriter opens a writer for the given request.
// A caller must close the writer after using it.
// Also, he must always check the Err() method.
//
// The method handles access control.
func (f *fileStore) OpenWriter(request *http.Request) io.WriteCloser {
	if f.hasErr() {
		return nil
	}
	f.assertHasWriteAccessForRequest(request)
	return f.openWriterAtPath(f.pathFromRequest(request))
}

// Delete will delete the file or directory pointed to by the request.
// A caller must always check the Err() method.
func (f *fileStore) Delete(request *http.Request) {
	if f.hasErr() {
		return
	}
	f.assertHasDeleteAccessForRequest(request)

	var p = f.pathFromRequest(request)
	if p.HasErr() {
		f.setErr(p.Err())
		return
	}

	var err = os.Remove(p.Path())
	f.setErr(err)
}
