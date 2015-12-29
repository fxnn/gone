package filestore

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/fxnn/gone/authenticator"
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

func sutNotAuthenticated(t *testing.T) store.Store {
	return New(getwd(t), authenticator.NewNeverAuthenticated())
}

func sutAuthenticated(t *testing.T) store.Store {
	return New(getwd(t), authenticator.NewAlwaysAuthenticated())
}

func requestGET(path string) (request *http.Request) {
	request, _ = http.NewRequest("GET", path, nil)
	return
}

func createTempSymlinkInCurrentwd(t *testing.T, target string) string {
	wd := getwd(t)
	symlinkName := path.Base(target) // let's use the same name
	symlink := path.Join(wd, symlinkName)
	if err := os.Symlink(target, symlink); err != nil {
		t.Fatalf("couldnt create symlink %s to %s: %s", symlink, target, err)
	}
	return symlinkName
}

func createTempWdInCurrentwd(t *testing.T, mode os.FileMode) string {
	wd := getwd(t)
	tempDirName := createTempDirInCurrentwd(t, mode)
	tempWd := path.Join(wd, tempDirName)
	if err := os.Chdir(tempWd); err != nil {
		t.Fatalf("couldnt change wd to %s: %s", tempWd, err)
	}
	return tempDirName
}

func createTempDirInCurrentwd(t *testing.T, mode os.FileMode) string {
	wd := getwd(t)
	tmpdir, err := ioutil.TempDir(wd, "gone_test_")
	if err != nil {
		t.Fatalf("couldnt create tempdir in %s: %s", wd, err)
	}
	err = os.Chmod(tmpdir, mode)
	if err != nil {
		t.Fatalf("couldnt chmod tempdir %s: %s", tmpdir, err)
	}
	return path.Base(tmpdir)
}

func createTempFileInCurrentwd(t *testing.T, mode os.FileMode) string {
	wd := getwd(t)
	tmpfile, err := ioutil.TempFile(wd, "gone_test_")
	if err != nil {
		t.Fatalf("couldnt create tempfile in %s: %s", wd, err)
	}
	info, err := tmpfile.Stat()
	if err != nil {
		t.Fatalf("couldnt stat tmpfile %s: %s", tmpfile, err)
	}
	err = tmpfile.Chmod(mode)
	if err != nil {
		t.Fatalf("couldnt chmod tmpfile %s: %s", info.Name(), err)
	}
	err = tmpfile.Close()
	if err != nil {
		t.Fatalf("couldn close tmpfile %s: %s", info.Name(), err)
	}
	return info.Name()
}

func removeTempSymlinkFromCurrentwd(t *testing.T, symlinkName string) {
	removeTempFileFromCurrentwd(t, symlinkName)
}

func removeTempWdFromCurrentwd(t *testing.T, tmpdir string) {
	newwd := path.Dir(getwd(t))
	if err := os.Chdir(newwd); err != nil {
		t.Fatalf("couldnt chdir to %s: %s", newwd, err)
	}
	removeTempDirFromCurrentwd(t, tmpdir)
}

func removeTempDirFromCurrentwd(t *testing.T, tmpdir string) {
	wd := getwd(t)
	tmpdirPath := path.Join(wd, tmpdir)
	err := os.Remove(tmpdirPath)
	if err != nil {
		t.Fatalf("couldnt remove tmpdir %s: %s", tmpdirPath, err)
	}
}

func removeTempFileFromCurrentwd(t *testing.T, tmpfile string) {
	wd := getwd(t)
	tmpfilePath := path.Join(wd, tmpfile)
	err := os.Remove(tmpfilePath)
	if err != nil {
		t.Fatalf("couldn remove tmpfile %s: %s", tmpfilePath, err)
	}
}

func getwd(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("couldnt get working directory: %s", err)
	}
	return wd
}

func closed(c io.Closer) {
	if c != nil {
		c.Close()
	}
}
