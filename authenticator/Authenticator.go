package authenticator

import (
	"net/http"
)

// Authenticator provides functions to deliver authentication information.
type Authenticator interface {
	// IsAuthenticated returns true iff the request indicates a properly
	// authenticated caller.
	IsAuthenticated(request *http.Request) bool

	// UserId returns a unique identifier of the user being currently logged
	// in.
	// If no user is logged in (therefore, if IsAuthenticated returns false),
	// UserId returns the empty string.
	UserId(request *http.Request) string
}

// HttpAuthenticator provides Authenticator's information and additionally
// allows the authentification of an user.
//
// Currently, it is tightly coupled to HTTP requests and even contains a UI
// implementation.
// Later, we should extract the UI part into the github.com/fxnn/gone/http
// package.
type HttpAuthenticator interface {
	Authenticator

	// AuthHandler can be part of the handler chain to read authentication
	// information from the session associated with the request, prior to
	// calling the delegate handler.
	AuthHandler(delegate http.Handler) http.Handler

	// ServeHTTP serves an authentication UI.
	ServeHTTP(writer http.ResponseWriter, request *http.Request)
}
