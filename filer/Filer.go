package filer

import (
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"

	"github.com/fxnn/gone/authenticator"
)

const fallbackMimeType = "application/octet-stream"

// Filer maps incoming HTTP requests to the file system.
type Filer struct {
	accessControl
}

// New initializes a zeroe'd instance ready to use.
func New(authenticator authenticator.Authenticator) *Filer {
	return &Filer{newAccessControl(authenticator)}
}

func (f *Filer) MimeTypeForRequest(request *http.Request) string {
	if f.err != nil {
		return ""
	}
	return f.mimeTypeForPath(f.pathFromRequest(request))
}

func (f *Filer) mimeTypeForPath(symlinkedPath string) string {
	// NOTE, that filename based mimetype detection needs symlinks resolution
	var p = f.evalSymlinks(symlinkedPath)
	if f.isDirectory(p) || f.err != nil {
		return fallbackMimeType
	}

	var ext = path.Ext(p)
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}

	var first512Bytes = f.first512BytesForPath(p)
	f.Err() // clear error flag, as DetectContentType always returns something

	result := http.DetectContentType(first512Bytes)
	return result
}

func (f *Filer) first512BytesForPath(p string) []byte {
	if f.err != nil {
		return nil
	}

	var readCloser = f.openReaderAtPath(p)
	if f.err != nil {
		return nil
	}
	var buf []byte = make([]byte, 512)
	var n int
	n, f.err = readCloser.Read(buf)
	readCloser.Close()

	return buf[:n]
}

// ReadString returns the requested content as string.
// A caller must always check the Err() method.
func (f *Filer) ReadString(request *http.Request) string {
	if f.err != nil {
		return ""
	}
	return f.readAllAndClose(f.OpenReader(request))
}

// WriteString writes the given content into a file pointed to by the request.
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
//
// The method handles access control.
func (f *Filer) OpenReader(request *http.Request) io.ReadCloser {
	if f.err != nil {
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
func (f *Filer) OpenWriter(request *http.Request) io.WriteCloser {
	if f.err != nil {
		return nil
	}
	f.assertHasWriteAccessForRequest(request)
	return f.openWriterAtPath(f.pathFromRequest(request))
}

func (f *Filer) openReaderAtPath(p string) (reader io.ReadCloser) {
	f.assertPathValidForAnyAccess(p)
	if f.err != nil {
		return nil
	}
	reader, err := os.Open(p)
	f.setErr(err)
	return
}

func (f *Filer) openWriterAtPath(p string) (writer io.WriteCloser) {
	f.assertPathValidForAnyAccess(p)
	if f.err != nil {
		return nil
	}
	writer, err := os.Create(p)
	f.setErr(err)
	return
}

// HtpasswdFilePath returns the path to the ".htpasswd" file in the content
// root, if one exists.
// Otherwise, it returns the empty string and sets the Err() value.
func (f *Filer) HtpasswdFilePath() string {
	htpasswdFilePath := path.Join(f.contentRootPath, ".htpasswd")
	f.assertPathExists(htpasswdFilePath)
	if f.err != nil {
		return ""
	}
	return htpasswdFilePath
}
