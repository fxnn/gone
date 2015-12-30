package editor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fxnn/gone/store"
)

func TestWriteSuccess(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = postRequest(t, "/someFile", "")
	var store = store.NewMockStore()
	var sut = New(store)

	request.PostForm.Set("content", "content")
	store.GivenWriteAccess()
	sut.serveWriter(response, request)

	assertResponseBody(t, response, "")
	assertResponseCode(t, response, 302)
	assertResponseHeader(t, response, "Location", "/someFile?edit")

}

func TestWriteUnauthorized(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = postRequest(t, "/someFile", "")
	var store = store.NewMockStore()
	var sut = New(store)

	sut.serveWriter(response, request)

	assertResponseBodyNotEmpty(t, response)
	assertResponseCode(t, response, 401)

}

func assertResponseBodyNotEmpty(t *testing.T, response *httptest.ResponseRecorder) {
	if response.Body.String() == "" {
		t.Fatalf("body expected to be non-empty, but is empty")
	}
}
func assertResponseBody(t *testing.T, response *httptest.ResponseRecorder, expected string) {
	if response.Body.String() != "" {
		t.Fatalf("body expected to be %s, but is %v", expected, response.Body.String())
	}
}

func assertResponseCode(t *testing.T, response *httptest.ResponseRecorder, expected int) {
	if response.Code != expected {
		t.Fatalf("code expected to be %d, but is %d", expected, response.Code)
	}
}

func assertResponseHeader(t *testing.T, response *httptest.ResponseRecorder, key string, expected string) {
	if response.Header().Get(key) != expected {
		t.Fatalf("%s header expected to be %s, but is %v", key, expected, response.Header().Get(key))
	}
}

func postRequest(t *testing.T, requestUrl string, content string) *http.Request {
	var request, err = http.NewRequest("POST", requestUrl, strings.NewReader(content))
	if err != nil {
		t.Fatalf("couldn't create http.Request: %v", err)
	}
	request.PostForm = url.Values{}
	return request
}
