package filestore

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/fxnn/gone/store"
)

// basicFiler implements basic algorithms and error handling needed when
// dealing with files.
// It basically wraps golang library functions for error handling.
type basicFiler struct {
	contentRoot string
	*errStore
}

func newBasicFiler(contentRoot string, errStore *errStore) *basicFiler {
	return &basicFiler{
		contentRoot: contentRoot,
		errStore:    errStore,
	}
}

func (f *basicFiler) stat(p string) (result os.FileInfo) {
	if f.hasErr() {
		return nil
	}

	if result, f.err = os.Stat(p); f.hasErr() {
		f.err = f.wrapErr(f.err)
		return nil
	}

	return result
}

func (f *basicFiler) pathFromRequest(request *http.Request) string {
	var p = f.guessExtension(f.normalizePath(path.Join(f.contentRoot, request.URL.Path)))
	if f.hasErr() {
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
	if f.hasErr() {
		return dir
	}
	var index = f.guessExtension(path.Join(dir, "index"))
	f.assertPathExists(index)
	if err := f.errAndClear(); err != nil {
		return dir
	}
	return index
}

// guessExtension tries to append the file extension, if missing.
// If the given path points to a valid file,
// simply returns the argument.
// Otherwise, it looks for all files in the
// directory beginning with the filename and a dot ("."), and returns the first
// match in alphabetic order.
// Err() will not be set.
func (f *basicFiler) guessExtension(p string) string {
	if f.hasErr() {
		return p
	}
	if f.assertPathExists(p); f.errAndClear() == nil {
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
	if f.hasErr() {
		return path
	}

	return f.cleanPath(f.absPath(path))
}

func (f *basicFiler) absPath(path string) (absPath string) {
	if f.hasErr() {
		return path
	}
	absPath, err := filepath.Abs(path)
	f.setErr(err)
	return
}

func (f *basicFiler) assertPathExists(path string) {
	if f.hasErr() {
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		f.setErr(store.NewPathNotFoundError(err.Error()))
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
	if f.hasErr() {
		return path
	}
	hardPath, err := filepath.EvalSymlinks(path)
	f.setErr(err)
	return
}

func (f *basicFiler) cleanPath(path string) string {
	if f.hasErr() {
		return path
	}
	return filepath.Clean(path)
}
