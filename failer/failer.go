package failer

import (
	"io"
	"net/http"
)

// The failer serves HTTP responses with error messages.
type failer struct {
	message string
	code    int
}

func newFailer(message string, code int) failer {
	return failer{message, code}
}

func (h failer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(h.code)
	io.WriteString(writer, h.message)
}
