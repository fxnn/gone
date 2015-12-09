package gopath

import "os"

// GoPath is an immutable object, representing a computational stage in
// path processing.
// It might either represent a path after zero or more successfully executed
// operations, or it might represent an errorneous state after an operation that
// returned an error.
type GoPath struct {
	path string
	err  error

	// A cached Stat() result, if available
	fileInfo os.FileInfo
}

var empty GoPath = GoPath{}

// Empty returns the empty GoPath.
func Empty() GoPath {
	return empty
}

// FromPath constructs a GoPath instance with the given path.
func FromPath(p string) GoPath {
	return GoPath{path: p}
}

// FromError constructs an errorneous GoPath instance.
func FromErr(err error) GoPath {
	return GoPath{err: err}
}

func (g GoPath) ClearErr() GoPath {
	return g.withErr(nil)
}

func (g GoPath) withPath(p string) GoPath {
	// NOTE, that we don't keep the fileInfo here, as another path could point
	// to another file
	return GoPath{path: p, err: g.err}
}

func (g GoPath) withErr(err error) GoPath {
	return GoPath{g.path, err, g.fileInfo}
}

func (g GoPath) withFileInfo(fileInfo os.FileInfo) GoPath {
	return GoPath{g.path, g.err, fileInfo}
}

type GoPathTransformer func(GoPath) GoPath

func (g GoPath) Do(transformer GoPathTransformer) GoPath {
	if g.HasErr() {
		return g
	}
	return transformer(g)
}
