package filer

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

// Maps incoming HTTP requests to the file system.
type accessControl struct {
	basicFiler
}

func (f *accessControl) HasWriteAccessForRequest(request *http.Request) bool {
	return f.hasWriteAccessToPath(f.pathFromRequest(request))
}

func (f *accessControl) assertHasWriteAccessToPath(p string) {
	if f.err == nil && !f.hasWriteAccessToPath(p) {
		f.setErr(NewAccessDeniedError(fmt.Sprintf("Access denied on %s", p)))
	}
}

func (f *accessControl) hasWriteAccessToPath(p string) bool {
	if f.err != nil {
		return false
	}
	info, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		// HINT: Inspect permissions of containing directory
		info, err = os.Stat(path.Dir(p))
	}
	if err != nil {
		f.setErr(err)
		return false
	}
	return f.hasWriteAccessForFileMode(info.Mode())
}

func (f *accessControl) hasWriteAccessForFileMode(mode os.FileMode) bool {
	// 0002 is the write permission write for others
	return mode&0002 != 0
}
