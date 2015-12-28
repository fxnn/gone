// Package authenticator allows to identify a user.
// That is, the package implements the logic behind login functionality
// as well as the session logic used to identify the user in all following
// requests.
//
// The implementation may be exchanged; so this package aims to allow using
// HTTP basic auth and OpenID at the same time.
package authenticator
