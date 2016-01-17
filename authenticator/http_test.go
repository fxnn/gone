package authenticator

import (
	"net/http"
	"testing"
)

func TestAuthAttemptUser(t *testing.T) {
	var sut = HttpBasicAuthenticator{}
	var req = &http.Request{Header: make(http.Header)}
	var actualUser string

	actualUser = sut.userAttemptingAuth(req)
	if actualUser != "" {
		t.Fatalf("Expected empty result, was '%v'", actualUser)
	}

	// from https://en.wikipedia.org/wiki/Basic_access_authentication
	req.Header.Add("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")

	actualUser = sut.userAttemptingAuth(req)
	if actualUser != "Aladdin" {
		t.Fatalf("Expected Aladdin, was '%v'", actualUser)
	}
}
