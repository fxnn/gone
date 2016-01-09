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
[![Coverage Status](https://coveralls.io/repos/fxnn/gone/badge.svg?branch=master&service=github)](https://coveralls.io/github/fxnn/gone?branch=master)


## Usage

Install the application and start it.

```console
$ go get github.com/fxnn/gone
$ gone
```

The current working directory will now be served on port `8080`.
For example, the file `test.md` in that working directory is now accessible as `http://localhost:8080/test.md`.

Append `?edit` to the URL to edit the content.
Append `?create` to the URL to create a non-existant file.

While no one will be able to access the outside of the working directory (e.g. by using `/../breakout`),
it _is_ possible to access symlinks to anywhere in the file system.

See `gone -help` for usage information and configuration options.


## Access Control

Gone uses the file system's access control features.
Of course, the Gone process can't read or write files it doesn't have a
permission to.
For example, if the Gone process is run by user `joe`, it won't be able to read
a file only user `ann` has read permission for (as with `rw-------`).

Likewise, an anonymous user being not logged in can't read or write files
through Gone, except those who have _world_ permissions.
For example, a file `rw-rw-r--` might be read by an anonymous user, but he
won't be able to change that file.
Also, in a directory `rwxrwxr-x`, only a user being logged in may create new files.

Users can login by appending `?login` to the URL.
The login information is configured in a good old `.htpasswd` file, placed in the working directory
of the Gone process.
Authenticated users can read and write all files that are readable
resp. writeable by the Gone process.

Note that there's a brute force blocker.
After each failed login attempt, the login request will be answered with an
increasing delay of up to 10 seconds.
The request delay is imposed per user, per IP address and globally.
The global delay, however, grows ten times slower than the other delays.

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


## Templates

Gone uses some Go templates for its UI.
The templates are shipped inside the executable, but you can use custom versions of them.
For general information on Go HTML templates, see the [html/template godoc](https://golang.org/pkg/html/template/).

With your web root as working directory, invoke `gone export-templates`.
It creates a new folder `.templates` which will never be delivered via HTTP.
You'll find all templates inside and can modify them.
If you (re)start Gone now, it will use the templates from that directory.

Note, that you can also supply a custom template path.
See `gone -help` for more information.


## Future

Some day, Gone might be able to
* use external programs to render files into HTML, which would allow you to display manpages or syntax-highlighted code right in the web browser.
* support OpenID authentication.
* respect each files / directories group attribute for access control, in combination with a `htgroup` file.
* incorporate Git as version control system.
* include a fulltext search.


## Development

If you want to modify sources in this project, you might find the following information helpful.


### Third party software

Please note that the project uses the vendoring tool https://github.com/kardianos/govendor.
Also, we use the standard go `vendor` folder, which means that all external projects are vendored and to be found in the `vendor` folder.
A list of projects and versions is managed under [vendor/vendor.json](vendor/vendor.json).
If you build with go-1.5, enable the [GO15VENDOREXPERIMENT](https://golang.org/s/go15vendor) flag.

Gone imports code from following projects:

* [abbot/go-http-auth](https://github.com/abbot/go-http-auth) for HTTP basic authentication
* [gorilla](https://github.com/gorilla), a great web toolkit for Go, used for sessions and cookies
* [russross/blackfriday](https://github.com/russross/blackfriday), a well-made markdown processor for Go
* [shurcooL/sanitized_anchor_name](https://github.com/shurcooL/sanitized_anchor_name) for making strings URL-compatible
* [golang.org/x/crypto](https://golang.org/x/crypto) for session-related cryptography
* [golang.org/x/net/context](https://golang.org/x/net/context) for request-scoped values
* [fxnn/gopath](https://github.com/fxnn/gopath) for easy handling of filesystem paths

Also, the following commands are used to build gone:

* [pierre/gotestcover](https://github.com/pierrre/gotestcover) to run tests with coverage analysis on multiple packages
* [mjibson/esc](https://github.com/mjibson/esc) for embedding files into the binary


### Architecture

                   +------+
                   | main |
                   +------+
                    |  | |
          +---------+  | +---------+
          v            v           v
      +-------+    +------+    +--------+
      | store |    | http |    | config |
      +-------+    +------+    +--------+
                   /   |  \
         +--------+    |   +--------+
         v             v            v
    +--------+    +--------+    +--------+
    | viewer |    | editor |    | router |
    +--------+    +--------+    +--------+

`main` just implements the startup logic and integrates all other top-level
components.
Depending on what `config` returns, a command is executed, which by default
starts up the web server.
From now on, we have to main parts.

On the one hand, there is the `store` that implements the whole storage.
Currently, the only usable storage engine is the filesystem.

On the other hand, there is the `http` package that serves HTTP requests using
different handlers.
The `router` component directs each request to the matching handler.
Handlers are implemented in the `viewer` and the `editor` package.
While the `editor` serves the editing UI, the `viewer` is responsible for 
serving whatever file is requested.

Other noteable packages are as follows.
* The `http/failer` package delivers error pages for HTTP requests.
* The `http/templates` package caches and renders the templates used for HTML
  output.
* The `resources` package encapsulates access to static resources, which are
  bundled with each `gone` distribution.

See the [Godoc](http://godoc.org/github.com/fxnn/gone) for more information.


## License (MIT)

Licensed under the MIT License, see [LICENSE](LICENSE) file for more information.
