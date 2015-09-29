package filer

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestOpenWriterSupportsCreatingFiles(t *testing.T) {
	tmpdir := createTempDirInCurrentwd(t, 0777)
	defer removeTempDirFromCurrentwd(t, tmpdir)

	sut := New()

	writeCloser := sut.OpenWriter(requestGET("/" + tmpdir + "/newFile"))
	if err := sut.Err(); err != nil {
		t.Fatalf("failed to open file for writing: %s", err)
	}

	writeCloser.Close()
}

func TestOpenWriterDeniesWhenWorldPermissionIsMissing(t *testing.T) {
	tmpfile := createTempFileInCurrentwd(t, 0770)
	defer removeTempFileFromCurrentwd(t, tmpfile)

	sut := New()

	writeCloser := sut.OpenWriter(requestGET("/" + tmpfile))
	err := sut.Err()
	if err == nil || !IsAccessDeniedError(err) {
		writeCloser.Close()
		t.Fatalf("expected AccessDeniedError on %s, but got %s", tmpfile, err)
	}
}

func requestGET(path string) (request *http.Request) {
	request, _ = http.NewRequest("GET", path, nil)
	return
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

func removeTempDirFromCurrentwd(t *testing.T, tmpdir string) {
	wd := getwd(t)
	tmpdirPath := path.Join(wd, tmpdir)
	err := os.RemoveAll(tmpdirPath)
	if err != nil {
		t.Fatalf("couldnt remove tmpdir %s: %s", tmpdirPath, err)
	}
}

func removeTempFileFromCurrentwd(t *testing.T, tmpfile string) {
	wd := getwd(t)
	tmpfilePath := path.Join(wd, tmpfile)
	err := os.RemoveAll(tmpfilePath)
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
