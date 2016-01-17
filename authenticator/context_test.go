package authenticator

import (
	"net/http"
	"testing"
)

func TestSetUserID(t *testing.T) {

	var sut = NewContextAuthenticator()
	var req = &http.Request{}
	var rsw = http.ResponseWriter(nil)

	sut.SetUserID(rsw, req, "test")
	if actual := sut.UserID(req); actual != "test" {
		t.Fatalf("After being set to '%v', userID was '%v'", "test", actual)
	}
	if !sut.IsAuthenticated(req) {
		t.Fatalf("After setting userID, user was not authenticated")
	}

}

func TestUserIDInitiallyZero(t *testing.T) {

	var sut = NewContextAuthenticator()
	var req = &http.Request{}

	if actual := sut.UserID(req); actual != "" {
		t.Fatalf("Expected userID to be initially zero, but was '%v'", actual)
	}
}

func TestInitiallyNotAuthenticated(t *testing.T) {

	var sut = NewContextAuthenticator()
	var req = &http.Request{}

	if sut.IsAuthenticated(req) {
		t.Fatalf("Initially authenticated")
	}

}
