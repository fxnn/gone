package gopath

import (
	"os"
	"path"
	"path/filepath"
	"strings"
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

// Glob calls filepath.Glob(string) and returns the matches and any possible
// error.
//
// See also GlobAny(), if you just want one GoPath value.
// Especially it returns g.Err(), if this GoPath is errorneous.
func (g GoPath) Glob() (matches []string, err error) {
	if g.HasErr() {
		return make([]string, 0), g.Err()
	}
	return filepath.Glob(g.Path())
}

// Components returns all path components in this gopath.
// That is, the path is split at each os.PathSeparator.
//
// If this gopath is errorneous, it returns the empty array.
func (g GoPath) Components() []string {
	if g.HasErr() {
		return make([]string, 0)
	}
	return strings.Split(g.Path(), string(os.PathSeparator))
}
