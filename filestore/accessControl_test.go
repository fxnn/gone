package filestore

import (
	"testing"

	"github.com/fxnn/gone/store"
)

// https://github.com/fxnn/gone/issues/7 No.2 positive
func TestCreateFileProceedsWithSupplementaryPermissions(t *testing.T) {
	tmpdir := createTempDirInCurrentwd(t, 0773) // World execute and write flags
	defer removeTempDirFromCurrentwd(t, tmpdir)

	sut := sutNotAuthenticated(t)

	writeCloser := sut.OpenWriter(requestGET("/" + tmpdir + "/newFile"))
	closed(writeCloser)
	if err := sut.Err(); err != nil {
		t.Fatalf("failed to open file for writing: %s", err)
	}
	removeTempFileFromCurrentwd(t, tmpdir+"/newFile")
}

// https://github.com/fxnn/gone/issues/7 No.2 negative write
func TestCreateFileDeniesWhenWorldWritePermissionIsMissing(t *testing.T) {
	tmpdir := createTempDirInCurrentwd(t, 0771) // World execute flag
	defer removeTempDirFromCurrentwd(t, tmpdir)

	sut := sutNotAuthenticated(t)

	writeCloser := sut.OpenWriter(requestGET("/" + tmpdir + "/newFile"))
	closed(writeCloser)
	if err := sut.Err(); err == nil || !store.IsAccessDeniedError(err) {
		t.Fatalf("expected AccessDeniedError on %s, but got %s", tmpdir+"/newFile", err)
	}
}

// https://github.com/fxnn/gone/issues/7 No.2 negative execute
func TestCreateFileDeniesWhenWorldExecutePermissionIsMissing(t *testing.T) {
	tmpdir := createTempDirInCurrentwd(t, 0772) // World write flag
	defer removeTempDirFromCurrentwd(t, tmpdir)

	sut := sutNotAuthenticated(t)

	writeCloser := sut.OpenWriter(requestGET("/" + tmpdir + "/newFile"))
	closed(writeCloser)
	if err := sut.Err(); err == nil || !store.IsAccessDeniedError(err) {
		t.Fatalf("expected AccessDeniedError on %s, but got %s", tmpdir+"/newFile", err)
	}
}

// https://github.com/fxnn/gone/issues/7 No.3 positive
func TestReadExistingFileInsideDirectoryProceedsWithSupplementaryPermissions(t *testing.T) {
	tmpdir := createTempWdInCurrentwd(t, 0771) // world execute flag
	defer removeTempWdFromCurrentwd(t, tmpdir)

	tmpfile := createTempFileInCurrentwd(t, 0774) // world read flag
	defer removeTempFileFromCurrentwd(t, tmpfile)

	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenReader(requestGET("/" + tmpfile))
	closed(readCloser)
	if err := sut.Err(); err != nil {
		t.Fatalf("failed to open file for reading: %s", err)
	}
}

// https://github.com/fxnn/gone/issues/7 No.3 negative (execute)
func TestReadExistingFileInsideDirectoryDeniesOnMissingExecutePermission(t *testing.T) {
	tmpdir := createTempWdInCurrentwd(t, 0770) // missing world execute flag
	defer removeTempWdFromCurrentwd(t, tmpdir)

	tmpfile := createTempFileInCurrentwd(t, 0774) // world read flag
	defer removeTempFileFromCurrentwd(t, tmpfile)

	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenReader(requestGET("/" + tmpdir + "/" + tmpfile))
	closed(readCloser)
	if err := sut.Err(); err == nil || !store.IsAccessDeniedError(err) {
		t.Fatalf("expected AccessDeniedError on %s, but got %s", tmpfile, err)
	}
}

// https://github.com/fxnn/gone/issues/7 No.3 negative (read)
func TestReadExistingFileInsideDirectoryDeniesOnMissingReadPermission(t *testing.T) {
	tmpdir := createTempWdInCurrentwd(t, 0771) // world execute flag
	defer removeTempWdFromCurrentwd(t, tmpdir)

	tmpfile := createTempFileInCurrentwd(t, 0772) // missing world read flag
	defer removeTempFileFromCurrentwd(t, tmpfile)

	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenReader(requestGET("/" + tmpfile))
	closed(readCloser)
	if err := sut.Err(); err == nil || !store.IsAccessDeniedError(err) {
		t.Fatalf("expected AccessDeniedError on %s, but got %s", tmpfile, err)
	}
}

// https://github.com/fxnn/gone/issues/7 No.4 positive
func TestWriteExistingFileInsideDirectoryProceedsWithSupplementaryPermissions(t *testing.T) {
	tmpdir := createTempWdInCurrentwd(t, 0771) // world execute flag
	defer removeTempWdFromCurrentwd(t, tmpdir)

	tmpfile := createTempFileInCurrentwd(t, 0772) // world write flag
	defer removeTempFileFromCurrentwd(t, tmpfile)

	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenWriter(requestGET("/" + tmpfile))
	closed(readCloser)
	if err := sut.Err(); err != nil {
		t.Fatalf("failed to open file for writing: %s", err)
	}
}

// https://github.com/fxnn/gone/issues/7 No.4 negative (execute)
func TestWriteExistingFileInsideDirectoryDeniesOnMissingExecutePermission(t *testing.T) {
	tmpdir := createTempWdInCurrentwd(t, 0770) // missing world execute flag
	defer removeTempWdFromCurrentwd(t, tmpdir)

	tmpfile := createTempFileInCurrentwd(t, 0772) // world write flag
	defer removeTempFileFromCurrentwd(t, tmpfile)

	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenWriter(requestGET("/" + tmpdir + "/" + tmpfile))
	closed(readCloser)
	if err := sut.Err(); err == nil || !store.IsAccessDeniedError(err) {
		t.Fatalf("expected AccessDeniedError on %s, but got %s", tmpfile, err)
	}
}

// https://github.com/fxnn/gone/issues/7 No.4 negative (write)
func TestWriteExistingFileInsideDirectoryDeniesOnMissingWritePermission(t *testing.T) {
	tmpdir := createTempWdInCurrentwd(t, 0771) // world execute flag
	defer removeTempWdFromCurrentwd(t, tmpdir)

	tmpfile := createTempFileInCurrentwd(t, 0770) // missing world write flag
	defer removeTempFileFromCurrentwd(t, tmpfile)

	sut := sutNotAuthenticated(t)

	readCloser := sut.OpenWriter(requestGET("/" + tmpfile))
	closed(readCloser)
	if err := sut.Err(); err == nil || !store.IsAccessDeniedError(err) {
		t.Fatalf("expected AccessDeniedError on %s, but got %s", tmpfile, err)
	}
}

func TestOpenReaderProceedsWhenAuthenticated(t *testing.T) {
	tmpfile := createTempFileInCurrentwd(t, 0770)
	defer removeTempFileFromCurrentwd(t, tmpfile)

	sut := sutAuthenticated(t)

	readCloser := sut.OpenReader(requestGET("/" + tmpfile))
	closed(readCloser)
	if err := sut.Err(); err != nil {
		t.Fatalf("failed to open %s for reading: %s", tmpfile, err)
	}
}

func TestAccessToParentDirDenied(t *testing.T) {
	tempFile := createTempFileInCurrentwd(t, 0777)
	defer removeTempFileFromCurrentwd(t, tempFile)

	tempWd := createTempWdInCurrentwd(t, 0777)
	defer removeTempWdFromCurrentwd(t, tempWd)

	sut := sutAuthenticated(t)

	readCloser := sut.OpenReader(requestGET("/../" + tempFile))
	closed(readCloser)
	if err := sut.Err(); err == nil {
		t.Fatalf("could open reader for parent dir of working directory %s", getwd(t))
	} else if !store.IsPathNotFoundError(err) {
		t.Fatalf("expected PathNotFoundError: %s", err)
	}
}

func TestAccessToSymlinkToParentDirAllowed(t *testing.T) {
	tempFile := createTempFileInCurrentwd(t, 0777)
	defer removeTempFileFromCurrentwd(t, tempFile)

	tempWdName := createTempWdInCurrentwd(t, 0777)
	defer removeTempWdFromCurrentwd(t, tempWdName)

	symlinkName := createTempSymlinkInCurrentwd(t, "../"+tempFile)
	defer removeTempSymlinkFromCurrentwd(t, symlinkName)

	sut := sutAuthenticated(t)

	readCloser := sut.OpenReader(requestGET("/" + symlinkName))
	closed(readCloser)
	if err := sut.Err(); err != nil {
		t.Fatalf("could open reader for symlink to file %s from wd %s: %s", tempFile, getwd(t), err)
	}
}
