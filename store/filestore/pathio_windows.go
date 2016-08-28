package filestore

import "github.com/fxnn/gopath"
import "os"
import "math/rand"
import "strconv"

func isPathWriteable(p gopath.GoPath) bool {
	if p.IsDirectory() {
		return isDirWriteable(p)
	}

	return isFileWriteable(p)
}

func isFileWriteable(p gopath.GoPath) bool {
	var closer, err = os.OpenFile(p.Path(), os.O_WRONLY, 0)
	if closer != nil {
		closer.Close()
	}

	return err == nil
}

func isDirWriteable(p gopath.GoPath) bool {
	var nonExistantFile = p.JoinPath("githubcom-fxnn-gone")
	for !nonExistantFile.HasErr() && nonExistantFile.IsExists() {
		nonExistantFile = nonExistantFile.Append("-" + strconv.Itoa(rand.Int()))
	}
	if nonExistantFile.HasErr() {
		return false
	}

	var closer, err = os.Create(nonExistantFile.Path())
	if closer != nil {
		closer.Close()
		os.Remove(nonExistantFile.Path())
	}

	return err == nil
}
