# gone

Gone is a wiki written in [Go](http://golang.org). It's

* KISS,
* Convention over Configuration and
* designed with Developers and Admins in mind.

It displays Markdown (as HTML), HTML and Plaintext straight from the filesystem.
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

While no one will be able to access the outside of the working directory (e.g. by using `/../breakout`),
it _is_ possible to access symlinks to anywhere in the file system.


## Access Control

Gone uses the file system's access control features.
Of course, the Gone process can't read or write files it doesn't have a permission to.
For example, if the Gone process is run by user `joe`, it won't be able to read a file only user `ann` has read permission for (as with `rw-------`).

Likewise, an anonymous user being not logged in can't read or write files through Gone, except those who have _world_ permissions.
For example, a file `rw-rw-r--` might be read by an anonymous user, but he won't be able to change that file.
Also, in a directory `rwxrwxr-x`, only a user being logged in may create new files.

Users can login by appending `?login` to the URL.
The login information is configured in a good old `.htpasswd` file, placed in the working directory
of the Gone process.
Authenticated users can read and write all files that are readable
resp. writeable by the Gone process.

### Security considerations

* Authentication information are submitted without encryption, so *use SSL*!
* Anyone may read *and write* files just by assigning world read/write permissions, so better
  `chmod -R o-rw *` if you want to keep your stuff secret!
* Gone uses the working directory for content delivery, so better use a start script which
  invokes `cd`!


## Index documents, file names

Calling a directory, Gone will look for a file named `index`.
Calling any file that does not exist (including `index`), Gone will try to look
for files with a extension appended and use the first one in alphabetic order.

So, the file `http://localhost:8080/test.md` could also be referenced as
`http://localhost:8080/test`, as long as no `test` file exists.
In the same way, an `index.md` file can be used as index document and will fulfill
the above requirements.

This mechanism is transparent to the user, no redirect will happen.


## Future

Some day, Gone might be able to
* use external programs to render files into HTML, which would allow you to display manpages or syntax-highlighted code right in the web browser.
* use customized templates, so that you can change the appereance of the wiki just be editing some Go HTML template files.
* support OpenID authentication.
* respect each files / directories group attribute for access control, in combination with a `htgroup` file.
* incorporate Git as version control system.
* include a fulltext search.


## Development

If you want to modify sources in this project, you might find the following information helpful.


### Third party software

Please note that the project uses the vendoring tool https://github.com/kardianos/govendor.
Consequently, all external projects are vendored and to be found in the `internal` folder.
A list of projects and versions is managed under [internal/vendor.json](internal/vendor.json)

Gone imports code from following projects:

* [abbot/go-http-auth](https://github.com/abbot/go-http-auth) for HTTP basic authentication
* [gorilla](https://github.com/gorilla), a great web toolkit for Go, used for sessions and cookies
* [russross/blackfriday](https://github.com/russross/blackfriday), a well-made markdown processor for Go
* [shurcooL/sanitized_anchor_name](https://github.com/shurcooL/sanitized_anchor_name) for making strings URL-compatible
* [golang.org/x/crypto](https://golang.org/x/crypto) for session-related cryptography
* [golang.org/x/net/context](https://golang.org/x/net/context) for request-scoped values
* [fxnn/gopath](https://github.com/fxnn/gopath) for easy handling of filesystem paths


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
* The `filestore` encapsulates mapping requests to the filesystem as well as reading
  and writing the files themselves.
* The `templates` package caches and renders the templates used for HTML output.
* The `resources` package encapsulates access to static resources, which are
  bundled with each `gone` distribution.
* The `failer` package delivers error pages for HTTP requests.

See the [Godoc](http://godoc.org/github.com/fxnn/gone) for more information.
