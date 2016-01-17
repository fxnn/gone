package authenticator

import (
	"net/http"
)

// Authenticator provides functions to deliver authentication information.
type Authenticator interface {
	// IsAuthenticated returns true iff the request indicates a properly
	// authenticated caller.
	IsAuthenticated(request *http.Request) bool

	// UserID returns a unique identifier of the user being currently logged
	// in.
	// If no user is logged in (therefore, if IsAuthenticated returns false),
	// UserID returns the empty string.
	UserID(request *http.Request) string

	// SetUserID sets the unique identifier of the user being currently logged
	// in.
	// Set this to the empty string to make no user being logged in.
	SetUserID(writer http.ResponseWriter, request *http.Request, userID string)
}

// HttpAuthenticator allows the authentification of an user over an HTTP
// protocol.
//
// Currently, it is tightly coupled to HTTP requests and even contains a UI
// implementation.
// Later, we should extract the UI part into the github.com/fxnn/gone/http
// package.
type HttpAuthenticator interface {
	// AuthHandler can be part of the handler chain to read authentication
	// information from the session associated with the request, prior to
	// calling the delegate handler.
	MiddlewareHandler(delegate http.Handler) http.Handler

	// LoginHandler provides a handler that serves the login UI.
	LoginHandler() http.Handler
}
