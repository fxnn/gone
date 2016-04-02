# gone

Gone is a wiki engine written in [Go](http://golang.org). It's

* KISS,
* Convention over Configuration and
* designed with Developers and Admins in mind.

With Gone, you can

* display Markdown, HTML and Plaintext straight from the filesystem.
* edit just any file that's made of text.
* have all this without setup, no database needed, not even the tinyest configuration.

So go get it!

[![Build Status](https://travis-ci.org/fxnn/gone.svg?branch=master)](https://travis-ci.org/fxnn/gone)
[![GoDoc](https://godoc.org/github.com/fxnn/gone?status.svg)](https://godoc.org/github.com/fxnn/gone)
[![Coverage Status](https://coveralls.io/repos/fxnn/gone/badge.svg?branch=master&service=github)](https://coveralls.io/github/fxnn/gone?branch=master)


## Usage

> *NOTE: This assumes that you have [Go installed](https://golang.org/doc/install).
> Binary distributions will follow.*

Install the application and start it.

```console
$ go get github.com/fxnn/gone
$ gone
```

The current working directory will now be served on port `8080`.

* *Display content.*
  The file `test.md` in that working directory is now accessible as `http://localhost:8080/test.md`.
  It's a [Markdown](https://en.wikipedia.org/wiki/Markdown) file, but Gone delivers a rendered webpage.
  Other files (text, HTML, PDF, ...) would simply be rendered as they are.
* *Editing just anything that's made of text.*
  In your browser, append `?edit` in the address bar.
  Gone now sends you a text editor, allowing you to edit your file.
  Your file doesn't exist yet? Use `?create` instead.
* *Customize everything.*
  Change how Gone looks.
  Call `gone export-templates`, and you will get the HTML, CSS and JavaScript behind Gone's frontend.
  Modify it as you like.

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

Some day, Gone might be
* extensible. Plugin in version control, renderers, compilers or anything you like. #29
* granting file access on a group level, using a `htgroup` file.
* searchable in full text.


## Development

If you want to modify sources in this project, you might find the following information helpful.


### Third party software

Please note that the project uses the vendoring tool https://github.com/kardianos/govendor.
Also, we use the standard go `vendor` folder, which means that all external projects are vendored and to be found in the `vendor` folder.
A list of projects and versions is managed under [vendor/vendor.json](vendor/vendor.json).
If you build with go-1.5, enable the [GO15VENDOREXPERIMENT](https://golang.org/s/go15vendor) flag.

Gone imports code from following projects:

* [abbot/go-http-auth](https://github.com/abbot/go-http-auth) for HTTP basic authentication
* [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify) for watching files
* [gorilla](https://github.com/gorilla), a great web toolkit for Go, used for sessions and cookies
* [russross/blackfriday](https://github.com/russross/blackfriday), a well-made markdown processor for Go
* [shurcooL/sanitized_anchor_name](https://github.com/shurcooL/sanitized_anchor_name) for making strings URL-compatible
* [golang.org/x/crypto](https://golang.org/x/crypto) for session-related cryptography
* [golang.org/x/net/context](https://golang.org/x/net/context) for request-scoped values
* [fxnn/gopath](https://github.com/fxnn/gopath) for easy handling of filesystem paths

Also, the following commands are used to build gone:

* [pierre/gotestcover](https://github.com/pierrre/gotestcover) to run tests with coverage analysis on multiple packages
* [mjibson/esc](https://github.com/mjibson/esc) for embedding files into the binary

Gone's frontend wouldn't be anything without

* [ajaxorg/ace](https://github.com/ajaxorg/ace), a great in-browser editor


## Architecture

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
