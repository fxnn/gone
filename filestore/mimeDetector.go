package filestore

import (
	"mime"
	"net/http"
	"path"

	"github.com/fxnn/gone/store"
)

type mimeDetector struct {
	*pathIO
	*basicFiler
	*errStore
}

func newMimeDetector(p *pathIO, f *basicFiler, s *errStore) *mimeDetector {
	return &mimeDetector{p, f, s}
}

func (m *mimeDetector) mimeTypeForPath(p string) string {
	if m.isDirectory(p) || m.hasErr() {
		return store.FallbackMimeType
	}

	var ext = path.Ext(p)
	if mimeType := mime.TypeByExtension(ext); mimeType != "" {
		return mimeType
	}

	var first512Bytes = m.first512BytesForPath(p)
	m.errAndClear() // clear error flag, as DetectContentType always returns something

	return http.DetectContentType(first512Bytes)
}

func (m *mimeDetector) first512BytesForPath(p string) []byte {
	if m.hasErr() {
		return nil
	}

	var readCloser = m.openReaderAtPath(p)
	if m.hasErr() {
		return nil
	}
	var buf []byte = make([]byte, 512)
	var n int
	n, m.err = readCloser.Read(buf)
	readCloser.Close()

	return buf[:n]
}
