package filestore

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/fxnn/gone/store"
)

type pathIO struct {
	*basicFiler
	*errStore
}

func newPathIO(f *basicFiler, s *errStore) *pathIO {
	return &pathIO{f, s}
}

func (i *pathIO) openReaderAtPath(p string) (reader io.ReadCloser) {
	i.assertPathValidForAnyAccess(p)
	if i.hasErr() {
		return nil
	}
	reader, err := os.Open(p)
	i.setErr(err)
	return
}

func (i *pathIO) openWriterAtPath(p string) (writer io.WriteCloser) {
	i.assertPathValidForAnyAccess(p)
	if i.hasErr() {
		return nil
	}
	writer, err := os.Create(p)
	i.setErr(err)
	return
}

func (i *pathIO) assertPathValidForAnyAccess(p string) {
	i.assertFileIsNotHidden(p)
	i.assertPathInsideContentRoot(p)
}

func (i *pathIO) assertFileIsNotHidden(p string) {
	if i.hasErr() {
		return
	}

	if strings.HasPrefix(path.Base(p), ".") {
		i.setErr(store.NewPathNotFoundError(fmt.Sprintf("%s is a hidden file and may not be displayed", p)))
	}
}

func (i *pathIO) assertPathInsideContentRoot(p string) {
	if i.hasErr() {
		return
	}

	var normalizedPath = i.normalizePath(p)

	if !i.hasErr() && !strings.HasPrefix(normalizedPath, i.contentRoot) {
		i.setErr(store.NewPathNotFoundError(
			fmt.Sprintf("%s is not inside content root %s", p, i.contentRoot),
		))
	}
}
