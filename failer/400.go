package failer

import (
	"net/http"
)

var (
	BadRequestHandler       = handler{"Oops, bad request", http.StatusBadRequest}
	NotFoundHandler         = handler{"Sorry, not found", http.StatusNotFound}
	MethodNotAllowedHandler = handler{"Oops, method not allowed", http.StatusMethodNotAllowed}
)

func ServeBadRequest(writer http.ResponseWriter, request *http.Request) {
	BadRequestHandler.ServeHTTP(writer, request)
}

func ServeNotFound(writer http.ResponseWriter, request *http.Request) {
	NotFoundHandler.ServeHTTP(writer, request)
}

func ServeMethodNotAllowed(writer http.ResponseWriter, request *http.Request) {
	MethodNotAllowedHandler.ServeHTTP(writer, request)
}
