package authenticator

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/abbot/go-http-auth"
	"github.com/fxnn/gone/authenticator/bruteblocker"
)

func TestMiddlewareHandler_noAuthentication(t *testing.T) {
	var requestAuth = newMockAuthenticator()
	var sessionAuth = newMockAuthenticator()
	var sut = sutWithRequestAuthAndSessionAuth(requestAuth, sessionAuth)
	var delegate = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {})
	var req = blankRequest()
	var rsp = newMockResponseWriter()

	sut.MiddlewareHandler(delegate).ServeHTTP(rsp, req)

	if requestAuth.IsAuthenticated(req) || sessionAuth.IsAuthenticated(req) {
		t.Fatalf("No authentication expected")
	}
}

func TestMiddlewareHandler_copiesUserId(t *testing.T) {
	var requestAuth = newMockAuthenticator()
	var sessionAuth = newMockAuthenticator()
	var sut = sutWithRequestAuthAndSessionAuth(requestAuth, sessionAuth)
	var delegate = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {})
	var req = blankRequest()
	var rsp = newMockResponseWriter()

	sessionAuth.SetUserID(rsp, req, "Aladdin")
	sut.MiddlewareHandler(delegate).ServeHTTP(rsp, req)

	if userID := requestAuth.UserID(req); userID != "Aladdin" {
		t.Fatalf("Expected Aladdin to be authenticated, but was '%v'", userID)
	}
}

func TestLoginHandler_rightCredentials(t *testing.T) {
	var requestAuth = newMockAuthenticator()
	var sessionAuth = newMockAuthenticator()
	var sut = sutWithBruteblockerRequestAuthSessionAuthAndBasicAuth(requestAuth, sessionAuth, aladdinsSecret)
	var req = requestWithAladdinUser()
	var rsp = newMockResponseWriter()

	sut.LoginHandler().ServeHTTP(rsp, req)

	if userID := requestAuth.UserID(req); userID != "Aladdin" {
		t.Fatalf("Expected request user to be Aladdin, but was '%v'", userID)
	}
	if userID := sessionAuth.UserID(req); userID != "Aladdin" {
		t.Fatalf("Expected session user to be Aladdin, but was '%v'", userID)
	}
	if rsp.status != http.StatusFound {
		t.Fatalf("Expected status authorized, was '%v'", rsp.status)
	}
}

func TestLoginHandler_wrongCredentials(t *testing.T) {
	var requestAuth = newMockAuthenticator()
	var sut = sutWithBruteblockerRequestAuthAndBasicAuth(requestAuth, noSecrets)
	var req = requestWithAladdinUser()
	var rsp = newMockResponseWriter()

	sut.LoginHandler().ServeHTTP(rsp, req)

	if rsp.status != http.StatusUnauthorized {
		t.Fatalf("Expected status unauthorized, was '%v'", rsp.status)
	}
	var userID = requestAuth.UserID(req)
	if userID != "" {
		t.Fatalf("Expected no authorization, but got user '%v'", userID)
	}
}

func TestLoginHandler_noLoginGiven(t *testing.T) {
	var sut = sutWithBasicAuth(noSecrets)
	var req = blankRequest()
	var rsp = newMockResponseWriter()

	sut.LoginHandler().ServeHTTP(rsp, req)

	if rsp.status != http.StatusUnauthorized {
		t.Fatalf("Expected status unauthorized, was '%v'", rsp.status)
	}
}

func TestAuthAttemptUser_withoutLogin(t *testing.T) {
	var sut = blankSut()
	var req = blankRequest()

	var actualUser = sut.userAttemptingAuth(req)

	if actualUser != "" {
		t.Fatalf("Expected empty result, was '%v'", actualUser)
	}
}

func TestAuthAttemptUser_successfulLogin(t *testing.T) {
	var sut = blankSut()
	var req = requestWithAladdinUser()

	var actualUser = sut.userAttemptingAuth(req)

	if actualUser != "Aladdin" {
		t.Fatalf("Expected Aladdin, was '%v'", actualUser)
	}
}

func sutWithRequestAuthAndSessionAuth(
	requestAuth Authenticator,
	sessionAuth Authenticator,
) (result HttpBasicAuthenticator) {
	result = blankSut()
	result.requestAuth = requestAuth
	result.sessionAuth = sessionAuth
	return
}

func sutWithBruteblockerRequestAuthSessionAuthAndBasicAuth(
	requestAuth Authenticator,
	sessionAuth Authenticator,
	secrets auth.SecretProvider,
) (result HttpBasicAuthenticator) {
	result = sutWithBruteblockerRequestAuthAndBasicAuth(requestAuth, secrets)
	result.sessionAuth = sessionAuth
	return
}

func sutWithBruteblockerRequestAuthAndBasicAuth(
	requestAuth Authenticator,
	secrets auth.SecretProvider,
) (result HttpBasicAuthenticator) {
	result = sutWithBasicAuth(secrets)
	result.requestAuth = requestAuth
	result.bruteBlocker = bruteblocker.New(0, 0, 0, 0, 0)
	return
}

func sutWithBasicAuth(secrets auth.SecretProvider) (result HttpBasicAuthenticator) {
	result = blankSut()
	result.basicAuth = auth.NewBasicAuthenticator("Realm", secrets)
	return
}

func blankSut() HttpBasicAuthenticator {
	return HttpBasicAuthenticator{}
}

func requestWithAladdinUser() (result *http.Request) {
	result = blankRequest()

	// from https://en.wikipedia.org/wiki/Basic_access_authentication
	result.Header.Add("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")

	return result
}

func blankRequest() *http.Request {
	return &http.Request{Header: make(http.Header), URL: &url.URL{}}
}
