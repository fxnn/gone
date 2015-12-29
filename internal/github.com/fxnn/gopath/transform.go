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
//
// Be warned: Stat() might cause an errorneous path to be returned, even in
// normal operation (e.g. file does not exist).
// An errorneous GoPath will have all operations being no-ops, so take care
// when using this function.
//
// Note, that Stat() is only useful for caching purposes.
// FileInfo() delivers the Stat() results even if Stat() was not called
// explicitly.
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
//
// If the path is already absolute, it returns the path itself.
// Otherwise, it returns an absolute representation of the path using the
// current working directory.
//
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
//
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

// GlobAny runs Glob() and selects the first match.
//
// If any error occurs, it returns an errorneous GoPath.
// Note, that -- according to https://godoc.org/path/filepath#Glob --
// this may only occur when the glob expression is not formatted
// correctly.
//
// If there is no match, an empty GoPath is returned.
func (g GoPath) GlobAny() GoPath {
	matches, err := g.Glob()
	if err != nil {
		return FromErr(err)
	}
	if len(matches) > 0 {
		return FromPath(matches[0])
	}
	return Empty()
}

// Rel returns the other (targpath) GoPath, expressed as path relative to this
// GoPath.
//
// 		var base = gopath.FromPath("/a")
//		var target = gopath.FromPath("/b/c")
//		var rel = base.Rel(target)
//
//		assert.Equal(rel.Path(), "../b/c")
//
// Note that this func follows the argument order of the filepath.Rel func,
// while the RelTo() func implements the reverse argument order.
func (g GoPath) Rel(targpath GoPath) GoPath {
	if g.HasErr() {
		return g
	}
	if targpath.HasErr() {
		return targpath
	}

	var result, err = filepath.Rel(g.Path(), targpath.Path())
	if err != nil {
		return targpath.withErr(err)
	}
	return targpath.withPath(result)
}

// RelTo returns this GoPath, expressed as path relative to the other (base)
// GoPath.
//
// 		var base = gopath.FromPath("/a")
//		var target = gopath.FromPath("/b/c")
//		var rel = target.RelTo(base)
//
//		assert.Equal(rel.Path(), "../b/c")
//
// Note that this func uses the inverse argument order of the filepath.Rel func,
// while the Rel() func implements the exact argument order.
func (g GoPath) RelTo(base GoPath) GoPath {
	if g.HasErr() {
		return g
	}

	return base.Rel(g)
}
