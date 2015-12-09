package gopath

import (
	"os"
	"path"
	"path/filepath"
)

func (g GoPath) Append(s string) GoPath {
	if g.HasErr() {
		return g
	}
	return g.withPath(g.Path() + s)
}

func (g GoPath) Join(other GoPath) GoPath {
	return g.JoinPath(other.Path())
}

func (g GoPath) JoinPath(p string) GoPath {
	if g.HasErr() {
		return g
	}
	return g.withPath(path.Join(g.Path(), p))
}

func (g GoPath) Dir() GoPath {
	if g.HasErr() {
		return g
	}
	return g.withPath(path.Dir(g.Path()))
}

// Stat calls os.Stat and caches the FileInfo result inside the returned
// GoPath.
// When the Stat call fails, an errorneous GoPath is returned.
// Stat always calls os.Stat, even if the GoPath already contains a FileInfo.
func (g GoPath) Stat() GoPath {
	if g.HasErr() {
		return g
	}

	if fileInfo, err := os.Stat(g.path); err != nil {
		return g.withErr(err)
	} else {
		return g.withFileInfo(fileInfo)
	}
}

// Abs calls filepath.Abs() on the path.
// If the path is already absolute, it returns the path itself.
// Otherwise, it returns an absolute representation of the path using the
// current working directory.
// If an error occurs, it returns an errorneous GoPath.
func (g GoPath) Abs() GoPath {
	if g.HasErr() {
		return g
	}

	if absPath, err := filepath.Abs(g.path); err != nil {
		return g.withErr(err)
	} else {
		return g.withPath(absPath)
	}
}

// EvalSymlinks calls filepath.EvalSymlinks().
// It evaluates any symlinks in the path.
// If the path is relative, the result might be relative, too.
// If an error occurs, it returns an errorneous GoPath.
func (g GoPath) EvalSymlinks() GoPath {
	if g.HasErr() {
		return g
	}

	if hardPath, err := filepath.EvalSymlinks(g.path); err != nil {
		return g.withErr(err)
	} else {
		return g.withPath(hardPath)
	}
}

// Clean calls filepath.Clean().
// It returns the shortest path equivalent to the given path.
// It might not return an errorneous GoPath, unless the given GoPath is already
// errorneous.
func (g GoPath) Clean() GoPath {
	if g.HasErr() {
		return g
	}

	return g.withPath(filepath.Clean(g.path))
}

func (g GoPath) GlobAny() GoPath {
	matches, err := g.Glob()
	if err != nil {
		return FromError(err)
	}
	if len(matches) > 0 {
		return FromPath(matches[0])
	}
	return Empty()
}
