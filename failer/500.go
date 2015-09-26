package failer

import (
	"net/http"
)

var (
	InternalServerErrorHandler = handler{
		"Oops, internal server error",
		http.StatusInternalServerError,
	}
)

func ServeInternalServerError(writer http.ResponseWriter, request *http.Request) {
	InternalServerErrorHandler.ServeHTTP(writer, request)
}
