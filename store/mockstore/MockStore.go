package mockstore

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/fxnn/gone/store"
)

// MockStore is only for testing and provides configurable answers.
// By default, it returns only zero values.
type MockStore struct {
	err          error
	readAccess   bool
	writeAccess  bool
	deleteAccess bool
	mimeType     string
	exists       bool
}

func New() *MockStore {
	return &MockStore{exists: true}
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

func (s *MockStore) GivenNotExists() {
	s.exists = false
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
	if !s.exists {
		s.err = store.NewPathNotFoundError("mocked PathNotFoundError")
	}
	return nil
}

func (s *MockStore) OpenWriter(request *http.Request) io.WriteCloser {
	return nil
}

func (s *MockStore) ReadString(request *http.Request) string {
	if !s.exists {
		s.err = store.NewPathNotFoundError("mocked PathNotFoundError")
	}
	return ""
}

func (s *MockStore) WriteString(request *http.Request, content string) {
}

func (s *MockStore) Delete(request *http.Request) {
	if !s.exists {
		s.err = store.NewPathNotFoundError("mocked PathNotFoundError")
	}
}

func (s *MockStore) FileSizeForRequest(request *http.Request) int64 {
	if !s.exists {
		s.err = store.NewPathNotFoundError("mocked PathNotFoundError")
	}
	return 0
}

func (s *MockStore) ModTimeForRequest(request *http.Request) time.Time {
	if !s.exists {
		s.err = store.NewPathNotFoundError("mocked PathNotFoundError")
	}
	return time.Now()
}

func (s *MockStore) GivenMimeType(mimeType string) {
	s.mimeType = mimeType
}

func (s *MockStore) MimeTypeForRequest(request *http.Request) string {
	if !s.exists {
		s.err = store.NewPathNotFoundError("mocked PathNotFoundError")
	}
	return s.mimeType
}

func (s *MockStore) Err() error {
	var result = s.err
	s.err = nil
	return result
}
