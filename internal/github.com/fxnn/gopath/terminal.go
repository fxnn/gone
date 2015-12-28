package gopath

import (
	"os"
	"path"
	"path/filepath"
)

func (g GoPath) String() string {
	return g.path
}

// Path returns the path represented by the given GoPath object.
// If that GoPath is errorneous, a zero value is returned.
func (g GoPath) Path() string {
	return g.path
}

func (g GoPath) Ext() string {
	return path.Ext(g.Path())
}

func (g GoPath) Base() string {
	return path.Base(g.Path())
}

// FileInfo calls Stat() on GoPath and returns the resulting FileInfo.
// When Stat() results in an error, or GoPath is already errorneous,
// the result is a zero interface.
// When the Stat() result is already cached inside the GoPath, Stat() is not
// called again.
func (g GoPath) FileInfo() os.FileInfo {
	if !g.HasErr() && g.fileInfo == nil {
		return g.Stat().FileInfo()
	}
	return g.fileInfo
}

// FileMode returns g.FileInfo().Mode(), or 0, when FileInfo() fails.
func (g GoPath) FileMode() os.FileMode {
	if info := g.FileInfo(); info != nil {
		return info.Mode()
	}
	return 0
}

func (g GoPath) Glob() (matches []string, err error) {
	return filepath.Glob(g.Path())
}
