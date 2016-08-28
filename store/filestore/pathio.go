package filestore

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/fxnn/gone/store"
	"github.com/fxnn/gopath"
)

// pathIO implements basic operations on paths
type pathIO struct {
	contentRoot gopath.GoPath
	*errStore
}

func newPathIO(contentRoot gopath.GoPath, s *errStore) *pathIO {
	var result = &pathIO{contentRoot, s}
	result.contentRoot = result.contentRoot.Do(result.normalizePath)
	return result
}

func (i *pathIO) openReaderAtPath(p gopath.GoPath) (reader io.ReadCloser) {
	i.assertPathValidForAnyAccess(p)
	if i.hasErr() {
		i.prependErr(fmt.Sprintf("cannot open reader for '%s'", p))
		return nil
	}

	reader, err := os.Open(p.Path())
	i.setErr(err)
	i.prependErr(fmt.Sprintf("couldn't open reader for '%s'", p))

	return
}

func (i *pathIO) openWriterAtPath(p gopath.GoPath) (writer io.WriteCloser) {
	i.assertPathValidForAnyAccess(p)
	i.assertPathValidForWriteAccess(p)
	if i.hasErr() {
		i.prependErr(fmt.Sprintf("cannot open writer for '%s'", p))
		return nil
	}

	writer, err := os.Create(p.Path())
	i.setErr(err)
	i.prependErr(fmt.Sprintf("couldn't open writer for '%s'", p))

	return
}

func (i *pathIO) assertPathExists(p gopath.GoPath) {
	i.syncedErrs(p.AssertExists())
	i.prependErr(fmt.Sprintf("required path %s does not exist", p))
}

// assertPathValidForWriteAccess sets the error flag when the path may not be
// opened for writing by this process.
func (i *pathIO) assertPathValidForWriteAccess(p gopath.GoPath) {
	if i.hasErr() {
		return
	}
	if p.HasErr() {
		i.setErr(p.Err())
		return
	}

	if p.IsExists() {
		if !p.IsRegular() || !isPathWriteable(p) {
			i.setErr(store.NewAccessDeniedError(fmt.Sprintf(
				"path '%s' with mode %s denotes no regular file or no writeable directory",
				p.Path(), p.FileMode())))
		}
	} else {
		var d = p.Dir()
		if !isPathWriteable(d) {
			i.setErr(store.NewAccessDeniedError(
				"parent directory of '" + p.Path() + "' is not writeable"))
		}
	}
}

// assertPathValidForAnyAccess sets the error flag when the path may not be
// accessed through this application in general.
// User-specific access permissions are NOT regarded here.
func (i *pathIO) assertPathValidForAnyAccess(p gopath.GoPath) {
	if p.HasErr() {
		i.syncedErrs(p)
	} else {
		i.assertFileIsNotHidden(p)
		i.assertPathInsideContentRoot(p)
	}
}

func (i *pathIO) assertFileIsNotHidden(p gopath.GoPath) {
	if i.hasErr() {
		return
	}
	if p.HasErr() {
		i.syncedErrs(p)
		return
	}

	if strings.HasPrefix(p.Base(), ".") {
		i.setErr(store.NewPathNotFoundError(fmt.Sprintf("%s is a hidden file and may not be displayed", p)))
	}

	// HINT: recursive call, ending at content root
	if i.isPathInsideContentRoot(p) {
		i.assertFileIsNotHidden(p.ToSlash().Dir())
	}
}

func (i *pathIO) assertPathInsideContentRoot(p gopath.GoPath) {
	if i.hasErr() {
		return
	}

	if !i.isPathInsideContentRoot(p) {
		i.setErr(store.NewPathNotFoundError(
			fmt.Sprintf("%s is not inside content root %s", p, i.contentRoot),
		))
	}
}

func (i *pathIO) isPathInsideContentRoot(p gopath.GoPath) bool {
	var normalizedPath = i.normalizePath(p)

	if !normalizedPath.HasErr() {
		return strings.HasPrefix(normalizedPath.Path(), i.contentRoot.Path())
	}

	return false
}

// pathFromRequest maps the request to the filesystem.
// It returns a GoPath that might be errorneous.
func (i *pathIO) pathFromRequest(request *http.Request) gopath.GoPath {
	if i.hasErr() {
		return gopath.FromErr(i.err)
	}

	var p = i.contentRoot.JoinPath(request.URL.Path).Do(i.normalizePath).Do(i.guessExtension)

	if !p.HasErr() && p.IsDirectory() {
		return i.indexForDirectory(p)
	}

	return i.syncedErrs(p.PrependErr("couldn't retrieve path from request"))
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
	var match = p.Append(".*").GlobAny().PrependErr("couldn't guess extension")
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
