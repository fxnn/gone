package editor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fxnn/gone/http/templates"
	"github.com/fxnn/gone/store"
	"github.com/fxnn/gone/store/mockstore"
)

func TestWriteSuccess(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = postRequest(t, "/someFile", "")
	var store = mockstore.New()
	var sut = createSut(store)

	request.PostForm.Set("content", "content")
	store.GivenWriteAccess()
	sut.ServeHTTP(response, request)

	assertResponseBody(t, response, "")
	assertResponseCode(t, response, http.StatusFound)
	assertResponseHeader(t, response, "Location", "/someFile?edit")

}

func TestWriteUnauthorized(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = postRequest(t, "/someFile", "")
	var store = mockstore.New()
	var sut = createSut(store)

	sut.ServeHTTP(response, request)

	assertResponseBodyNotEmpty(t, response)
	assertResponseCode(t, response, http.StatusUnauthorized)

}

func TestCreateUISuccess(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = getRequest(t, "/someFile?create")
	var store = mockstore.New()
	var sut = createSut(store)

	store.GivenWriteAccess()
	store.GivenNotExists()
	sut.ServeHTTP(response, request)

	assertResponseBodyNotEmpty(t, response)
	assertResponseCode(t, response, http.StatusOK)

}

func TestCreateUIMissingWritePermission(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = getRequest(t, "/someFile?create")
	var store = mockstore.New()
	var sut = createSut(store)

	store.GivenNotExists()
	sut.ServeHTTP(response, request)

	assertResponseBodyNotEmpty(t, response)
	assertResponseCode(t, response, http.StatusUnauthorized)

}

func TestCreateUIAlreadyExists(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = getRequest(t, "/someFile?create")
	var store = mockstore.New()
	var sut = createSut(store)

	store.GivenWriteAccess()
	store.GivenMimeType("text/plain")
	sut.ServeHTTP(response, request)

	assertResponseBodyNotEmpty(t, response)
	assertResponseCode(t, response, http.StatusConflict)

}

func TestEditUISuccess(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = getRequest(t, "/someFile?edit")
	var store = mockstore.New()
	var sut = createSut(store)

	store.GivenReadAccess()
	store.GivenWriteAccess()
	store.GivenMimeType("text/plain")
	sut.ServeHTTP(response, request)

	assertResponseBodyNotEmpty(t, response)
	assertResponseCode(t, response, http.StatusOK)

}

func TestEditUIMissingReadPermission(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = getRequest(t, "/someFile?edit")
	var store = mockstore.New()
	var sut = createSut(store)

	store.GivenWriteAccess()
	store.GivenMimeType("text/plain")
	sut.ServeHTTP(response, request)

	assertResponseBodyNotEmpty(t, response)
	assertResponseCode(t, response, http.StatusUnauthorized)

}

func TestEditUIMissingWritePermission(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = getRequest(t, "/someFile?edit")
	var store = mockstore.New()
	var sut = createSut(store)

	store.GivenReadAccess()
	store.GivenMimeType("text/plain")
	sut.ServeHTTP(response, request)

	assertResponseBodyNotEmpty(t, response)
	assertResponseCode(t, response, http.StatusUnauthorized)

}

func TestEditUINotExists(t *testing.T) {

	var response = httptest.NewRecorder()
	var request = getRequest(t, "/someFile?edit")
	var store = mockstore.New()
	var sut = createSut(store)

	store.GivenNotExists()
	store.GivenReadAccess()
	store.GivenWriteAccess()
	sut.ServeHTTP(response, request)

	assertResponseBodyNotEmpty(t, response)
	assertResponseCode(t, response, http.StatusNotFound)

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

	// HINT: PostForm initialisieren, um POST-Parameter mocken zu können
	request.PostForm = url.Values{}

	return request
}

func getRequest(t *testing.T, requestUrl string) *http.Request {
	var request, err = http.NewRequest("GET", requestUrl, strings.NewReader(""))
	if err != nil {
		t.Fatalf("couldn't create http.Request: %v", err)
	}

	// HINT: GET-Parameter auswerten
	request.ParseForm()

	return request
}

func createSut(s store.Store) *Editor {
	var l = templates.NewStaticLoader()
	return New(l, s)
}
