package filestore

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fxnn/gone/store"
	"github.com/fxnn/gone/internal/github.com/fxnn/gopath"
)

type pathIO struct {
	*basicFiler
	*errStore
}

func newPathIO(f *basicFiler, s *errStore) *pathIO {
	return &pathIO{f, s}
}

func (i *pathIO) openReaderAtPath(p gopath.GoPath) (reader io.ReadCloser) {
	if p.HasErr() {
		return nil
	}

	i.assertPathValidForAnyAccess(p)
	if i.hasErr() {
		return nil
	}

	reader, err := os.Open(p.Path())
	i.setErr(err)

	return
}

func (i *pathIO) openWriterAtPath(p gopath.GoPath) (writer io.WriteCloser) {
	if p.HasErr() {
		return nil
	}

	i.assertPathValidForAnyAccess(p)
	if i.hasErr() {
		return nil
	}

	writer, err := os.Create(p.Path())
	i.setErr(err)
	return
}

func (i *pathIO) assertPathValidForAnyAccess(p gopath.GoPath) {
	i.assertFileIsNotHidden(p)
	i.assertPathInsideContentRoot(p)
}

func (i *pathIO) assertFileIsNotHidden(p gopath.GoPath) {
	if i.hasErr() || p.HasErr() {
		return
	}

	if strings.HasPrefix(p.Base(), ".") {
		i.setErr(store.NewPathNotFoundError(fmt.Sprintf("%s is a hidden file and may not be displayed", p)))
	}
}

func (i *pathIO) assertPathInsideContentRoot(p gopath.GoPath) {
	if i.hasErr() || p.HasErr() {
		return
	}

	var normalizedPath = i.normalizePath(p)

	if !p.HasErr() && !strings.HasPrefix(normalizedPath.Path(), i.contentRoot.Path()) {
		i.setErr(store.NewPathNotFoundError(
			fmt.Sprintf("%s is not inside content root %s", p, i.contentRoot),
		))
	}
}
