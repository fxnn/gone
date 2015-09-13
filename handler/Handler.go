package handler

import (
	"io"
	"log"
	"net/http"
	"os"
)

// Handles HTTP requests to the Gone wiki.
type Handler struct{}

// Initializes a zeroe'd instance ready to use.
func NewHandler() *Handler {
	return &Handler{}
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
	// TODO: Prohibit requests to ".."
	var path = "." + request.URL.Path
	var file, err = os.Open(path)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err.Error())
		h.serveNotFound(writer, request)
		return
	}

	h.serveFromReader(file, writer, request)
	file.Close()
}

func (h *Handler) serveNotFound(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusNotFound)
	io.WriteString(writer, "Oops, file not found")
}

func (h *Handler) serveFromReader(reader io.Reader, writer http.ResponseWriter, request *http.Request) {
	var written, err = io.Copy(writer, reader)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err.Error())
		h.serveInternalError(writer, request)
		return
	}

	log.Printf("%s %s: wrote %d bytes", request.Method, request.URL, written)
}

func (h *Handler) serveInternalError(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusInternalServerError)
	io.WriteString(writer, "Oops, internal server error")
}
