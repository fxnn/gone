package authenticator

import "net/http"

type mockResponseWriter struct {
	header        http.Header
	headerWritten bool
	status        int
	bytesWritten  int
}

func newMockResponseWriter() *mockResponseWriter {
	return &mockResponseWriter{header: make(http.Header)}
}

func (w *mockResponseWriter) Header() http.Header {
	return w.header
}

func (w *mockResponseWriter) Write(content []byte) (written int, err error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return
}

func (w *mockResponseWriter) WriteHeader(status int) {
	w.headerWritten = true
	w.status = status
}
