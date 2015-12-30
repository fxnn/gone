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
