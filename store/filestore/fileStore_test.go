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

func TestReadHiddenFileReturnsPathNotFoundError(t *testing.T) {
	skipOnWindows(t)

	tempFile := createPrefixedTempFileInCurrentwd(t, 0777, ".")
	defer removeTempFileFromCurrentwd(t, tempFile)

	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenReader(requestGET("/" + tempFile))
	closed(readCloser)
	if err := sut.Err(); err == nil {
		t.Fatalf("expected error, but got nil")
	} else if !store.IsPathNotFoundError(err) {
		t.Fatalf("expected PathNotFoundError, but got %v", err)
	}
}

// https://github.com/fxnn/gone/issues/15
func TestCreateFileInHiddenDirReturnsPathNotFoundError(t *testing.T) {
	tempDir := createPrefixedTempDirInCurrentwd(t, 0777, ".")
	defer removeTempDirFromCurrentwd(t, tempDir)

	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenWriter(requestGET("/" + tempDir + "/newFile"))
	closed(readCloser)
	if err := sut.Err(); err == nil {
		t.Fatalf("expected error, but got nil")
	} else if !store.IsPathNotFoundError(err) {
		t.Fatalf("expected PathNotFoundError, but got %v", err)
	}
}

// https://github.com/fxnn/gone/issues/15
func TestReadInHiddenDirReturnsPathNotFoundError(t *testing.T) {
	skipOnWindows(t)

	tempDir := createPrefixedTempWdInCurrentwd(t, 0777, ".")
	defer removeTempWdFromCurrentwd(t, tempDir)

	tempFile := createTempFileInCurrentwd(t, 0777)
	defer removeTempFileFromCurrentwd(t, tempFile)

	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenReader(requestGET("/" + tempFile))
	closed(readCloser)
	if err := sut.Err(); err == nil {
		t.Fatalf("expected error, but got nil")
	} else if !store.IsPathNotFoundError(err) {
		t.Fatalf("expected PathNotFoundError, but got %v", err)
	}
}
