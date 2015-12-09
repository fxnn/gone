// gopath package implements an alternative API to Go's path processing
// libraries.
// What's special to gopath is its ability to write path processing in a fluent
// way while doing error handling later.
// This is possible by using the immutable object GoPath, which either
// represents a path, or an error.
// All operations on an errorneous GoPath are no-ops, therefore the first error
// ever occured in a chain of operations will remain visible.
package gopath
