# gone

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

[![Build Status](https://travis-ci.org/fxnn/gone.svg?branch=master)](https://travis-ci.org/fxnn/gone)
[![GoDoc](https://godoc.org/github.com/fxnn/gone?status.svg)](https://godoc.org/github.com/fxnn/gone)


## Usage

Simply start the application.
The current working directory will now be served on port `8080`.
Append `?edit` to the URL to edit the content.
Append `?create` to the URL to create a non-existant file.


## Architecture

                    +------+
                    | main |
                    +------+
                    /   |  \
          +--------+    |   +--------+
          v             v            v
     +--------+    +--------+    +--------+
     | viewer |    | editor |    | router |
     +--------+    +--------+    +--------+
       |            /   |
       |   +-------+    |
       v   v            v
    +-------+    +-----------+
    | filer |    | templates |
    +-------+    +-----------+

`main` implements the HTTP Server component, using different handlers to serve
requests.
The `router` component directs each request to the matching handler.
Handlers are implemented in the `viewer` and the `editor` package.
While the `editor` serves the editing UI, the `viewer` is responsible for 
serving whatever file is requested.

Both use a set of backend packages.
* The `filer` encapsulate reading and writing files from the filesystem.
* The `templates` package caches and renders the templates used for HTML output.
* The `resources` package encapsulates access to static resources, which are
  bundled with each `gone` distribution.
* The `failer` package delivers error pages for HTTP requests.

