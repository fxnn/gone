package filer

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
)

func TestOpenWriterSupportsCreatingFiles(t *testing.T) {
	tmpdir := createTempDirInCurrentwd(t)
	defer removeTempDirFromCurrentwd(t, tmpdir)

	request, _ := http.NewRequest("GET", "/"+tmpdir+"/newFile", nil)
	sut := New()

	writeCloser := sut.OpenWriter(request)
	if err := sut.Err(); err != nil {
		t.Fatalf("failed to open file for writing: %s", err)
	}

	writeCloser.Close()
}

func createTempDirInCurrentwd(t *testing.T) string {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("couldnt get working directory: %s", err)
	}
	tmpdir, err := ioutil.TempDir(wd, "gone_test_")
	if err != nil {
		t.Fatalf("couldnt create tempdir in %s: %s", wd, err)
	}
	return path.Base(tmpdir)
}

func removeTempDirFromCurrentwd(t *testing.T, tmpdir string) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("couldnt get working directory: %s", err)
	}
	tmpdirPath := path.Join(wd, tmpdir)
	err = os.RemoveAll(tmpdirPath)
	if err != nil {
		t.Fatalf("couldnt remove tmpdir %s: %s", tmpdirPath, err)
	}
}
