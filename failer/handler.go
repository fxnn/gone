package failer

import (
	"io"
	"net/http"
)

type handler struct {
	message string
	code    int
}

func (h handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(h.code)
	io.WriteString(writer, h.message)
}
