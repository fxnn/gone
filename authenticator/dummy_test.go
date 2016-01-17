package authenticator

import (
	"net/http"
	"testing"
)

func TestAlwaysAuthenticated(t *testing.T) {

	var request, _ = http.NewRequest("GET", "", nil)
	var sut = NewAlwaysAuthenticated()

	if sut.UserID(request) != "" {
		t.Fatalf("UserID is not '', but '%v'", sut.UserID(request))
	}
	if !sut.IsAuthenticated(request) {
		t.Fatalf("UserID '' not returned as authenticated")
	}

}

func TestNeverAuthenticated(t *testing.T) {

	var request, _ = http.NewRequest("GET", "", nil)
	var sut = NewNeverAuthenticated()

	if sut.UserID(request) != "" {
		t.Fatalf("UserID is not '', but '%v'", sut.UserID(request))
	}
	if sut.IsAuthenticated(request) {
		t.Fatalf("UserID '' returned as authenticated")
	}

}
