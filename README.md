# gone

Gone is a wiki written in [Go](http://golang.org). It's

* KISS,
* Convention over Configuration and
* designed with Developers and Admins in mind.

It displays Markdown, HTML and Plaintext straight from the filesystem.
It allows you to edit just anything that has MIME type `text/*`.
It uses the filesystem's access control as far as possible.

[![Build Status](https://travis-ci.org/fxnn/gone.svg?branch=master)](https://travis-ci.org/fxnn/gone)
[![GoDoc](https://godoc.org/github.com/fxnn/gone?status.svg)](https://godoc.org/github.com/fxnn/gone)


## Usage

Simply start the application.
The current working directory will now be served on port `8080`.
For example, the file `test.md` in that working directory is now accessible as `http://localhost:8080/test.md`.

Append `?edit` to the URL to edit the content.
Append `?create` to the URL to create a non-existant file.


## Access Control

Gone uses the file system's access control features.
Of course, the Gone process can't read or write files it doesn't have a permission to.
For example, if the Gone process is run by user `joe`, it won't be able to read a file only user `ann` has read permission for (as with `rw-r-----`).

Likewise, an anonymous user being not logged in can't read or write files through Gone, except those who have _world_ permissions.
For example, a file `rw-r--r--` might be read by an anonymous user, but he won't be able to change that file.
Also, in a directory `rwxr-xr-x`, only a user being logged in may create new files.

### Currently unimplemented access control features

Users can login.
The login information is configured in a good old `htpasswd` file.

By default, authenticated users can read and write all files that are readable
resp. writeable by the Gone process.


## Future

Some day, Gone might be able to
* use external programs to render files into HTML, which would allow you to display manpages or syntax-highlighted code right in the web browser.
* use customized templates, so that you can change the appereance of the wiki just be editing some Go HTML template files.
* support OpenID authentication.
* respect each files / directories group attribute for access control, in combination with a `htgroup` file.
* incorporate Git as version control system.


## Development

If you want to modify sources in this project, please note that the project uses the vendoring tool http://github.com/kardianos/govendor.


### Architecture

                    +------+
                    | main |
                    +------+
                    /   |  \
          +--------+    |   +--------+
          v             v            v
     +--------+    +--------+    +--------+
     | viewer |    | editor |    | router |
     +--------+    +--------+    +--------+

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

See the [Godoc](http://godoc.org/github.com/fxnn/gone) for more information.
