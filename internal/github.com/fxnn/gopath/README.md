# gopath
An alternative Go API for operating on paths without the need for error handling between each step.

[![Build Status](https://travis-ci.org/fxnn/gopath.svg?branch=master)](https://travis-ci.org/fxnn/gopath)
[![GoDoc](https://godoc.org/github.com/fxnn/gopath?status.svg)](https://godoc.org/github.com/fxnn/gopath)

## Usage

Instead of doing a lot of error handling between each step ...

```go
var err error
var p = "/my/path/to/somewhere"

if p, err = filepath.Abs(p); err != nil {
  // handle error
}
if p, err = filepath.EvalSymlinks(p); err != nil {
  // handle error
}
if p, err = filepath.Rel(p, "other/path"); err != nil {
  // handle error
}

// go on
```

...gopath allows you to write the operations first, and do the error checking later.

```go
var p = gopath.FromPath("/my/path/to/somewhere").Abs().EvalSymlinks().Rel("other/path")

if p.HasErr() {
  // handle error
}

// go on
```

The idea: the GoPath-object encapsulates a path _and_ an error object.
Now, as soon as an error occured, each operation turns to a no-op.
That's why it suffices to check for errors in the end.

## License (MIT)

Licensed under the MIT License, see [LICENSE](LICENSE) file for more information.
