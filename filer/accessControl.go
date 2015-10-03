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

func (a *accessControl) HasWriteAccessForRequest(request *http.Request) bool {
	return a.hasWriteAccessToPath(a.pathFromRequest(request))
}

func (a *accessControl) HasReadAccessForRequest(request *http.Request) bool {
	return a.hasReadAccessToPath(a.pathFromRequest(request))
}

func (a *accessControl) assertHasWriteAccessToPath(p string) {
	if a.err == nil && !a.hasWriteAccessToPath(p) {
		a.setErr(NewAccessDeniedError(fmt.Sprintf("Access denied on %s", p)))
	}
}

func (a *accessControl) assertHasReadAccessToPath(p string) {
	if a.err == nil && !a.hasReadAccessToPath(p) {
		a.setErr(NewAccessDeniedError(fmt.Sprintf("Access denied on %s", p)))
	}
}

func (a *accessControl) hasWriteAccessToPath(p string) bool {
	if a.err != nil {
		return false
	}
	return a.hasWriteAccessForFileMode(a.relevantFileModeForPath(p))
}

func (a *accessControl) hasReadAccessToPath(p string) bool {
	if a.err != nil {
		return false
	}
	return a.hasReadAccessForFileMode(a.relevantFileModeForPath(p))
}

func (a *accessControl) hasWriteAccessForFileMode(mode os.FileMode) bool {
	// 0002 is the write permission for others
	return mode&0002 != 0
}

func (a *accessControl) hasReadAccessForFileMode(mode os.FileMode) bool {
	// 0004 is the read permission for others
	return mode&0004 != 0
}

// getRelevantFileModeForPath returns the FileMode for the given file or, when
// the file does not exist, its containing directory.
func (a *accessControl) relevantFileModeForPath(p string) os.FileMode {
	if a.err != nil {
		return 0
	}
	info, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		// HINT: Inspect permissions of containing directory
		info, err = os.Stat(path.Dir(p))
	}
	if err != nil {
		a.setErr(err)
		return 0
	}
	return info.Mode()
}
