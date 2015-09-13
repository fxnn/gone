package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	path := "." + request.URL.Path
	ok, err := isPathInsideWorkingDirectory(path)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err.Error())
		h.serveInternalServerError(writer, request)
		return
	}
	if !ok {
		log.Printf("%s %s: Not inside working directory", request.Method, request.URL)
		h.serveNotFound(writer, request)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		log.Printf("%s %s: %s", request.Method, request.URL, err.Error())
		h.serveNotFound(writer, request)
		return
	}

	h.serveFromReader(file, writer, request)
	file.Close()
}

func isPathInsideWorkingDirectory(path string) (bool, error) {
	normalizedPath, err := normalizePath(path)
	if err != nil {
		return false, fmt.Errorf("checking %s inside wd: %s", path, err)
	}

	wdPath, err := os.Getwd()
	if err != nil {
		return false, fmt.Errorf("checking %s inside wd: %s", path, err)
	}
	normalizedWdPath, err := normalizePath(wdPath)
	if err != nil {
		return false, fmt.Errorf("checking %s inside wd: %s", path, err)
	}

	return strings.HasPrefix(normalizedPath, normalizedWdPath), nil
}

func normalizePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path, fmt.Errorf("building abs path of %s: %s", path, err)
	}

	// TODO: Check whether existis

	hardPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return path, fmt.Errorf("removing symlinks from %s: %s", absPath, err)
	}

	// HINT: Remove .. and ., remove trailing slash
	cleanPath := filepath.Clean(hardPath)

	return cleanPath, nil
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
