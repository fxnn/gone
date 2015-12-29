// Package gopath implements an alternative API to Go's path processing
// libraries.
// It implements the recommondations from the article
// https://blog.golang.org/errors-are-values.
// For even more information, also have a look at the article
// http://blog.golang.org/error-handling-and-go.
//
// Idea
//
// What's special to gopath is its ability to write code that processes paths
// in a fluent way, while doing the error handling later.
// This is possible by using the immutable object gopath.GoPath, which either
// represents a path, or an error.
//
// 		var p1 = gopath.FromPath("/existing/file").Stat()
//		assert.Nil(p1.Err())
//
//		var p2 = gopath.FromPath("/doesnt/exist").Stat()
//		assert.NotNil(p2.Err())
//
// All operations on an errorneous GoPath are no-ops, therefore the first error
// ever occured in a chain of operations will remain visible.
// This way, you can work with GoPaths like you would with any other object.
//
//		var p = gopath.FromPath("/my/path/to/somewhere").Abs().EvalSymlinks().Rel(otherPath)
//
//		if p.HasErr() {
//			// handle error
//		}
//
//		// go on
//
package gopath
