package gopath

import "os"
import "fmt"

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

var empty = GoPath{}

// Empty returns the empty GoPath.
func Empty() GoPath {
	return empty
}

// FromPath constructs a GoPath instance with the given path.
func FromPath(p string) GoPath {
	return GoPath{path: p}
}

// FromErr constructs an errorneous GoPath instance.
func FromErr(err error) GoPath {
	return GoPath{err: err}
}

// ClearErr returns a GoPath instance with the same fields as this instance,
// except for the err field being nil.
func (g GoPath) ClearErr() GoPath {
	return g.withErr(nil)
}

// IsEmpty returns true exactly iff this GoPath contains the empty path:
//
//      assert.True(gopath.FromPath("").IsEmpty())
//      assert.False(gopath.FromPath("/some/empty/file").IsEmpty())
//
// Note, that this does not check the file size, contents or anything like this.
func (g GoPath) IsEmpty() bool {
	return g.path == ""
}

// PrependErr modifies this GoPath by prefixing its error text with the given
// string, followed by a colon and a space.
// When this GoPath doesn't have an err set, it is returned unchanged.
func (g GoPath) PrependErr(prefix string) GoPath {
	if g.HasErr() {
		return g.withErr(fmt.Errorf("%s: %s", prefix, g.Err().Error()))
	}

	return g
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

// Transformer is a func that transforms one GoPath instance into another.
type Transformer func(GoPath) GoPath

// Do executes the transformer on this GoPath.
// It therefore is an extension point. For example, you could define a
// transformer like the following:
//
//      func normalizePath(p gopath.GoPath) gopath.GoPath {
//          return p.Abs().Clean()
//      }
//
// Now, you can invoke it using Do:
//
//      var p = gopath.FromPath("some/path").Do(normalizePath)
//
func (g GoPath) Do(transformer Transformer) GoPath {
	if g.HasErr() {
		return g
	}
	return transformer(g)
}
