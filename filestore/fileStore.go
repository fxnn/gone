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
	*ioUtil
	*pathIO
	*basicFiler
	*errStore
	*mimeDetector
	*accessControl
}

// New initializes a zeroe'd instance ready to use.
func New(contentRoot string, authenticator authenticator.Authenticator) store.Store {
	var s = newErrStore()
	var f = newBasicFiler(contentRoot, s)
	var i = newIOUtil(s)
	var p = newPathIO(f, s)
	var m = newMimeDetector(p, f, s)
	var a = newAccessControl(authenticator, s, f)
	return &fileStore{i, p, f, s, m, a}
}

func (f *fileStore) Err() error {
	return f.errAndClear()
}

func (f *fileStore) MimeTypeForRequest(request *http.Request) string {
	if f.hasErr() {
		return ""
	}
	return f.mimeTypeForPath(f.evalSymlinks(f.pathFromRequest(request)))
}

// FileSizeForRequest returns the size of the underlying file in bytes, if any,
// or sets the Err() value.
func (f *fileStore) FileSizeForRequest(request *http.Request) int64 {
	p := f.pathFromRequest(request)
	if f.hasErr() {
		return -1
	}

	var info os.FileInfo
	if info = f.stat(p); f.hasErr() {
		return -1
	}

	return info.Size()
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
