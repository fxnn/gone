# gone

[![Build Status](https://travis-ci.org/fxnn/gone.svg?branch=master)](https://travis-ci.org/fxnn/gone)
[![GoDoc](https://godoc.org/github.com/fxnn/gone?status.svg)](https://godoc.org/github.com/fxnn/gone)

Gone is a wiki written in go. It's

* KISS,
* Convention over Configuration and
* designed with sysadmins in mind.

It might someday display plain text, markdown, godoc, manpages etc.
Files of unknown type are provided for download.
Plaintext files might be edited.

For permissions, the file system's access control is used as far as possible.
You can define user groups that must have read/write permission for anonymous/authorized access.
Authentication is configured in a good old `htpasswd` file.
Additional groups for authorized users might be supplied in a `htgroup` file.


## Architecture

            +--------+
            | server |
            +--------+
             /      \
            v        v
     +--------+    +--------+
     | viewer |    | editor |
     +--------+    +--------+
       |            /    |
       |   +-------+     |
       v   v             v
    +-------+    +-----------+
    | filer |    | templates |
    +-------+    +-----------+

The `server` is the HTTP Server component, using different handlers to serve
requests.
Handlers are implemented in the `viewer` and the `editor` package.
While the `editor` serves the editing UI, the `viewer` is responsible for 
serving whatever file is requested.

Both use a set of backend packages.
On the one hand this is `filer`, which encapsulate reading and writing files
from the filesystem.
On the other hand there is the `templates` package, which caches and renders
the templates used for HTML output.

