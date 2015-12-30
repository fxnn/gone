package filestore

import (
	"testing"

	"github.com/fxnn/gone/store"
)

func TestReadNonExistantFileReturnsPathNotFoundError(t *testing.T) {
	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenReader(requestGET("/thisFileShouldNeverExist"))
	closed(readCloser)
	if err := sut.Err(); err == nil {
		t.Fatalf("expected error, but got nil")
	} else if !store.IsPathNotFoundError(err) {
		t.Fatalf("expected PathNotFoundError, but got %v", err)
	}
}
