package filer

import (
	"io"
	"log"
	"net/http"
)

// Handles HTTP requests to the Gone wiki.
type Handler struct {
	filer Filer
}

// Initializes a zeroe'd instance ready to use.
func NewHandler() *Handler {
	return &Handler{NewFiler()}
}

func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		h.serveNonGET(writer, request)
		return
	}

	h.serveGET(writer, request)
}

func (h *Handler) serveNonGET(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusMethodNotAllowed)
	io.WriteString(writer, "Oops, method not allowed")
}

func (h *Handler) serveGET(writer http.ResponseWriter, request *http.Request) {
	var readCloser, err = h.filer.OpenReader(request)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err.Error())
		if isPathNotFoundError(err) {
			h.serveNotFound(writer, request)
		} else {
			h.serveInternalServerError(writer, request)
		}
		return
	}

	h.serveFromReader(readCloser, writer, request)
	readCloser.Close()
}

func (h *Handler) serveNotFound(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	io.WriteString(writer, "Oops, file not found")
}

func (h *Handler) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	written, err := io.Copy(writer, reader)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err.Error())
		h.serveInternalServerError(writer, request)
		return
	}

	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, written)
}

func (h *Handler) serveInternalServerError(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)
	io.WriteString(writer, "Oops, internal server error")
}
