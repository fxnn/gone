package filer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Maps incoming HTTP requests to the file system.
type Filer struct {
	err error
}

// Initializes a zeroe'd instance ready to use.
func New() Filer {
	return Filer{}
}

// Never forget to check for errors
func (f *Filer) Err() error {
	return f.err
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

func (f *Filer) assertPathInsideWorkingDirectory(path string) {
	if f.err != nil {
		return
	}

	var normalizedPath = f.normalizePath(path)
	var wdPath = f.normalizePath(f.workingDirectory())

	if f.err == nil && !strings.HasPrefix(normalizedPath, wdPath) {
		f.err = NewPathNotFoundError(fmt.Sprintf("%s is not inside working directory", path))
	}
}

func (f *Filer) normalizePath(path string) string {
	if f.err != nil {
		return path
	}

	var result string

	result = f.absPath(path)

	f.assertPathExists(result)
	result = f.evalSymlinks(result)

	// HINT: Remove .. and ., remove trailing slash
	return f.cleanPath(result)
}

func (f *Filer) absPath(path string) (absPath string) {
	if f.err != nil {
		return path
	}
	absPath, f.err = filepath.Abs(path)
	return
}

func (f *Filer) assertPathExists(path string) {
	if f.err != nil {
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f.err = NewPathNotFoundError(err.Error())
	}
}

func (f *Filer) evalSymlinks(path string) (hardPath string) {
	if f.err != nil {
		return path
	}
	hardPath, f.err = filepath.EvalSymlinks(path)
	return
}

func (f *Filer) cleanPath(path string) string {
	if f.err != nil {
		return path
	}
	return filepath.Clean(path)
}

func (f *Filer) workingDirectory() (wd string) {
	if f.err != nil {
		return ""
	}
	wd, f.err = os.Getwd()
	return
}
