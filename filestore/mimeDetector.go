package filestore

import (
	"mime"
	"net/http"

	"github.com/fxnn/gopath"
	"github.com/fxnn/gone/store"
)

type mimeDetector struct {
	*pathIO
	*errStore
}

func newMimeDetector(p *pathIO, s *errStore) *mimeDetector {
	return &mimeDetector{p, s}
}

func (m *mimeDetector) mimeTypeForPath(p gopath.GoPath) string {
	p = p.EvalSymlinks()
	if p.IsDirectory() || p.HasErr() {
		return store.FallbackMimeType
	}

	var ext = p.Ext()
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}

	var first512Bytes = m.first512BytesForPath(p)
	m.errAndClear() // clear error flag, as DetectContentType always returns something

	return http.DetectContentType(first512Bytes)
}

func (m *mimeDetector) first512BytesForPath(p gopath.GoPath) []byte {
	if p.HasErr() {
		return nil
	}

	var readCloser = m.openReaderAtPath(p)
	if m.hasErr() {
		return nil
	}
	var buf = make([]byte, 512)
	var n int
	n, m.err = readCloser.Read(buf)
	readCloser.Close()

	return buf[:n]
}
