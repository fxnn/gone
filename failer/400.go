package failer

import (
	"net/http"
)

var (
	BadRequestHandler       = newFailer("Oops, bad request", http.StatusBadRequest)
	UnauthorizedHandler     = newFailer("Oops, unauthorized", http.StatusUnauthorized)
	NotFoundHandler         = newFailer("Sorry, not found", http.StatusNotFound)
	MethodNotAllowedHandler = newFailer("Oops, method not allowed", http.StatusMethodNotAllowed)
	ConflictHandler         = newFailer("Sorry, there's a conflict", http.StatusConflict)
)

func ServeBadRequest(writer http.ResponseWriter, request *http.Request) {
	BadRequestHandler.ServeHTTP(writer, request)
}

func ServeUnauthorized(writer http.ResponseWriter, request *http.Request) {
	UnauthorizedHandler.ServeHTTP(writer, request)
}

func ServeNotFound(writer http.ResponseWriter, request *http.Request) {
	NotFoundHandler.ServeHTTP(writer, request)
}

func ServeMethodNotAllowed(writer http.ResponseWriter, request *http.Request) {
	MethodNotAllowedHandler.ServeHTTP(writer, request)
}

func ServeConflict(writer http.ResponseWriter, request *http.Request) {
	ConflictHandler.ServeHTTP(writer, request)
}
