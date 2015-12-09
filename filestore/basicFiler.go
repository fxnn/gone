package filestore

import (
	"net/http"

	"github.com/fxnn/gone/internal/github.com/fxnn/gopath"
)

// basicFiler implements basic algorithms and error handling needed when
// dealing with files.
// It basically wraps golang library functions for error handling.
type basicFiler struct {
	contentRoot gopath.GoPath
	*errStore
}

func newBasicFiler(contentRoot string, errStore *errStore) *basicFiler {
	return &basicFiler{
		contentRoot: gopath.FromPath(contentRoot),
		errStore:    errStore,
	}
}

func (f *basicFiler) setErrAndReturnPath(p gopath.GoPath) string {
	if !f.hasErr() {
		f.setErr(p.Err())
	}
	return p.Path()
}

func (f *basicFiler) pathFromRequest(request *http.Request) gopath.GoPath {
	var p = f.contentRoot.JoinPath(request.URL.Path).Do(f.normalizePath).Do(f.guessExtension)

	//var p = f.guessExtension(f.normalizePath(path.Join(f.contentRoot, request.URL.Path)))
	if !p.HasErr() && p.IsDirectory() {
		return f.indexForDirectory(p)
	}

	return p
}

// indexForDirectory finds the index document inside the given directory.
// On success, it returns the path to the index document, otherwise it simply
// returns the given path.
func (f *basicFiler) indexForDirectory(dir gopath.GoPath) gopath.GoPath {
	if f.hasErr() {
		return dir
	}
	var index = dir.JoinPath("index").Do(f.guessExtension).AssertExists()
	if index.HasErr() {
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
func (f *basicFiler) guessExtension(p gopath.GoPath) gopath.GoPath {
	var match = p.Append(".*").GlobAny()
	if match.Path() != "" {
		return match
	}
	return p
}

// normalizePath builds an absolute path and cleans it from ".." and ".", but
// doesn't resolve symlinks
func (f *basicFiler) normalizePath(p gopath.GoPath) gopath.GoPath {
	return p.Abs().Clean()
}
