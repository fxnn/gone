package filestore

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fxnn/gopath"
	"github.com/fxnn/gone/store"
)

// pathIO implements basic operations on paths
type pathIO struct {
	contentRoot gopath.GoPath
	*errStore
}

func newPathIO(contentRoot string, s *errStore) *pathIO {
	var result = &pathIO{gopath.FromPath(contentRoot), s}
	result.contentRoot = result.contentRoot.Do(result.normalizePath)
	return result
}

func (i *pathIO) openReaderAtPath(p gopath.GoPath) (reader io.ReadCloser) {
	i.assertPathValidForAnyAccess(p)
	if i.hasErr() {
		return nil
	}

	reader, err := os.Open(p.Path())
	i.setErr(err)

	return
}

func (i *pathIO) openWriterAtPath(p gopath.GoPath) (writer io.WriteCloser) {
	i.assertPathValidForAnyAccess(p)
	if i.hasErr() {
		return nil
	}

	writer, err := os.Create(p.Path())
	i.setErr(err)
	return
}

// assertPathValidForAnyAccess sets the error flag when the path may not be
// accessed through this application in general.
// User-specific access permissions are NOT regarded here.
func (i *pathIO) assertPathValidForAnyAccess(p gopath.GoPath) {
	if p.HasErr() {
		i.setErr(p.Err())
	} else {
		i.assertFileIsNotHidden(p)
		i.assertPathInsideContentRoot(p)
	}
}

func (i *pathIO) assertFileIsNotHidden(p gopath.GoPath) {
	if i.hasErr() {
		return
	}

	if strings.HasPrefix(p.Base(), ".") {
		i.setErr(store.NewPathNotFoundError(fmt.Sprintf("%s is a hidden file and may not be displayed", p)))
	}
}

func (i *pathIO) assertPathInsideContentRoot(p gopath.GoPath) {
	if i.hasErr() {
		return
	}

	var normalizedPath = i.normalizePath(p)

	if !p.HasErr() && !strings.HasPrefix(normalizedPath.Path(), i.contentRoot.Path()) {
		i.setErr(store.NewPathNotFoundError(
			fmt.Sprintf("%s is not inside content root %s", p, i.contentRoot),
		))
	}
}

// pathFromRequest maps the request to the filesystem.
// It returns a GoPath that might be errorneous.
func (i *pathIO) pathFromRequest(request *http.Request) gopath.GoPath {
	var p = i.contentRoot.JoinPath(request.URL.Path).Do(i.normalizePath).Do(i.guessExtension)

	if !p.HasErr() && p.IsDirectory() {
		return i.indexForDirectory(p)
	}

	return i.syncedErrs(p)
}

// indexForDirectory finds the index document inside the given directory.
// On success, it returns the path to the index document, otherwise it simply
// returns the given path.
//
// Doesn't set the Err() value.
func (i *pathIO) indexForDirectory(dir gopath.GoPath) gopath.GoPath {
	var index = dir.JoinPath("index").Do(i.guessExtension).AssertExists()
	if index.HasErr() {
		return dir
	}

	return index
}

// guessExtension tries to append the file extension, if missing.
// If the given path points to a valid file, simply returns the argument.
// Otherwise, it looks for all files in the directory beginning with the
// filename and a dot ("."), and returns the first match in alphabetic order.
func (i *pathIO) guessExtension(p gopath.GoPath) gopath.GoPath {
	var match = p.Append(".*").GlobAny()
	if match.Path() != "" {
		return i.syncedErrs(match)
	}
	return p
}

// normalizePath builds an absolute path and cleans it from ".." and ".", but
// doesn't resolve symlinks
func (i *pathIO) normalizePath(p gopath.GoPath) gopath.GoPath {
	return p.Abs().Clean()
}

// pathComponentsTo returns a list of all path components between the content
// root directory and the given file.
//
// For example, a file "dir/file.ext" inside the content root will return both
// "dir" and "file.ext" as components.
func (i *pathIO) pathComponentsTo(p gopath.GoPath) []string {
	return i.contentRoot.Rel(p).Components()
}
