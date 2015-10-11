package filer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// basicFiler implements basic algorithms and error handling needed when
// dealing with files.
// It basically wraps golang library functions for error handling.
type basicFiler struct {
	contentRootPath string
	err             error
}

func newBasicFiler() basicFiler {
	return basicFiler{
		contentRootPath: "",
		err:             nil,
	}
}

// SetContentRootPath changes the path to the directory all the content is in.
func (f *basicFiler) SetContentRootPath(contentRootPath string) {
	f.contentRootPath = f.normalizePath(contentRootPath)
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

// WrapErr wraps f.err to a filer-specific error, if possible
func wrapErr(err error) error {
	if os.IsNotExist(err) {
		if pathError, ok := err.(*os.PathError); ok {
			return NewPathNotFoundError("path not found: " + pathError.Path)
		}
		return NewPathNotFoundError(fmt.Sprintf("path not found: %s", err))
	}
	return err
}

// FileSizeForRequest returns the size of the underlying file in bytes, if any,
// or sets the Err() value.
func (f *basicFiler) FileSizeForRequest(request *http.Request) int64 {
	p := f.pathFromRequest(request)
	if f.err != nil {
		return -1
	}

	var info os.FileInfo
	if info, f.err = os.Stat(p); f.err != nil {
		return -1
	}

	return info.Size()
}

func (f *basicFiler) pathFromRequest(request *http.Request) string {
	var p = f.guessExtension(f.normalizePath(path.Join(f.contentRootPath, request.URL.Path)))
	if f.err != nil {
		return p
	}
	if f.isDirectory(p) {
		return f.indexForDirectory(p)
	}
	return p
}

// indexForDirectory finds the index document inside the given directory.
// On success, it returns the path to the index document, otherwise it simply
// returns the given path.
func (f *basicFiler) indexForDirectory(dir string) string {
	if f.err != nil {
		return dir
	}
	var index = f.guessExtension(path.Join(dir, "index"))
	f.assertPathExists(index)
	if err := f.Err(); err != nil {
		return dir
	}
	return index
}

func (f *basicFiler) assertPathValidForAnyAccess(p string) {
	f.assertFileIsNotHidden(p)
	f.assertPathInsideContentRoot(p)
}

func (f *basicFiler) assertFileIsNotHidden(p string) {
	if f.err != nil {
		return
	}

	if strings.HasPrefix(path.Base(p), ".") {
		f.setErr(NewPathNotFoundError(fmt.Sprintf("%s is a hidden file and may not be displayed", p)))
	}
}

func (f *basicFiler) assertPathInsideContentRoot(p string) {
	if f.err != nil {
		return
	}

	var normalizedPath = f.normalizePath(p)

	if f.err == nil && !strings.HasPrefix(normalizedPath, f.contentRootPath) {
		f.setErr(NewPathNotFoundError(
			fmt.Sprintf("%s is not inside content root %s", p, f.contentRootPath),
		))
	}
}

// guessExtension tries to append the file extension, if missing.
// If the given path points to a valid file,
// simply returns the argument.
// Otherwise, it looks for all files in the
// directory beginning with the filename and a dot ("."), and returns the first
// match in alphabetic order.
// Err() will not be set.
func (f *basicFiler) guessExtension(p string) string {
	if f.err != nil {
		return p
	}
	if f.assertPathExists(p); f.Err() == nil {
		// don't apply for existing files
		return p
	}
	if matches, err := filepath.Glob(p + ".*"); err == nil && len(matches) > 0 {
		return matches[0]
	} else if err != nil {
		log.Printf("guessExtension for %s: %s", p, err)
	}
	return p
}

// normalizePath builds an absolute path and cleans it from ".." and ".", but
// doesn't resolve symlinks
func (f *basicFiler) normalizePath(path string) string {
	if f.err != nil {
		return path
	}

	return f.cleanPath(f.absPath(path))
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

// isDirectory returns true iff the path points to a directory. Err() will
// never be set.
func (f *basicFiler) isDirectory(path string) bool {
	if info, err := os.Stat(path); err != nil {
		return false
	} else {
		return info.IsDir()
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
