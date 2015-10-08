package filer

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// The basicFiler implements basic algorithms and error handling needed when
// dealing with files.
// It basically wraps golang library functions for error handling.
type basicFiler struct {
	err error
}

func newBasicFiler() basicFiler {
	return basicFiler{nil}
}

// Never forget to check for errors.
// One call to this function resets the error state.
func (f *basicFiler) Err() error {
	var result = f.err
	f.err = nil
	return result
}

func (f *basicFiler) setErr(err error) {
	f.err = wrapErr(err)
}

// Wraps f.err to a filer-specific error, if possible
func wrapErr(err error) error {
	if os.IsNotExist(err) {
		if pathError, ok := err.(*os.PathError); ok {
			return NewPathNotFoundError("path not found: " + pathError.Path)
		}
		return NewPathNotFoundError(fmt.Sprintf("path not found: %s", err))
	}
	return err
}

func (f *basicFiler) pathFromRequest(request *http.Request) string {
	var p = "." + request.URL.Path
	f.assertPathInsideWorkingDirectory(p)
	return p
}

func (f *basicFiler) assertPathInsideWorkingDirectory(p string) {
	if f.err != nil {
		return
	}

	var normalizedPath = f.normalizePath(p)
	var wdPath = f.normalizePath(f.workingDirectory())

	if f.err == nil && !strings.HasPrefix(normalizedPath, wdPath) {
		f.setErr(NewPathNotFoundError(fmt.Sprintf("%s is not inside working directory", p)))
	} else if f.err != nil {
		var oldErr = f.err
		f.err = nil
		f.assertPathInsideWorkingDirectory(path.Dir(p))
		if f.err != nil {
			f.err = oldErr
		}
	}
}

// Builds an absolute path and cleans it from ".." and ".", but doesn't resolve
// symlinks
func (f *basicFiler) normalizePath(path string) string {
	if f.err != nil {
		return path
	}

	var result string

	result = f.absPath(path)
	f.assertPathExists(result)

	// HINT: Remove .. and ., remove trailing slash
	return f.cleanPath(result)
}

func (f *basicFiler) absPath(path string) (absPath string) {
	if f.err != nil {
		return path
	}
	absPath, err := filepath.Abs(path)
	f.setErr(err)
	return
}

func (f *basicFiler) assertPathExists(path string) {
	if f.err != nil {
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f.setErr(NewPathNotFoundError(err.Error()))
	}
}

func (f *basicFiler) evalSymlinks(path string) (hardPath string) {
	if f.err != nil {
		return path
	}
	hardPath, err := filepath.EvalSymlinks(path)
	f.setErr(err)
	return
}

func (f *basicFiler) cleanPath(path string) string {
	if f.err != nil {
		return path
	}
	return filepath.Clean(path)
}

func (f *basicFiler) workingDirectory() (wd string) {
	if f.err != nil {
		return ""
	}
	wd, err := os.Getwd()
	f.setErr(err)
	return
}
