package store

import (
	"errors"
	"io"
	"net/http"
)

// MockStore is only for testing and provides configurable answers.
// By default, it returns only zero values.
type MockStore struct {
	err          error
	readAccess   bool
	writeAccess  bool
	deleteAccess bool
}

func NewMockStore() *MockStore {
	return &MockStore{}
}

func (s *MockStore) GivenNoErr() {
	s.GivenErr(nil)
}

func (s *MockStore) GivenSomeErr() {
	s.GivenErr(errors.New("mock error"))
}

func (s *MockStore) GivenErr(err error) {
	s.err = err
}

func (s *MockStore) GivenReadAccess() {
	s.readAccess = true
}

func (s *MockStore) HasReadAccessForRequest(request *http.Request) bool {
	return s.readAccess
}

func (s *MockStore) GivenWriteAccess() {
	s.writeAccess = true
}

func (s *MockStore) HasWriteAccessForRequest(request *http.Request) bool {
	return s.writeAccess
}

func (s *MockStore) GivenDeleteAccess() {
	s.deleteAccess = true
}

func (s *MockStore) HasDeleteAccessForRequest(request *http.Request) bool {
	return s.deleteAccess
}

func (s *MockStore) OpenReader(request *http.Request) io.ReadCloser {
	return nil
}

func (s *MockStore) OpenWriter(request *http.Request) io.WriteCloser {
	return nil
}

func (s *MockStore) ReadString(request *http.Request) string {
	return ""
}

func (s *MockStore) WriteString(request *http.Request, content string) {
}

func (s *MockStore) Delete(request *http.Request) {
}

func (s *MockStore) FileSizeForRequest(request *http.Request) int64 {
	return 0
}

func (s *MockStore) MimeTypeForRequest(request *http.Request) string {
	return ""
}

func (s *MockStore) Err() error {
	return s.err
}
