package failer

import (
	"net/http"
)

var (
	BadRequestHandler       = newFailer("Oops, bad request", http.StatusBadRequest)
	NotFoundHandler         = newFailer("Sorry, not found", http.StatusNotFound)
	MethodNotAllowedHandler = newFailer("Oops, method not allowed", http.StatusMethodNotAllowed)
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
