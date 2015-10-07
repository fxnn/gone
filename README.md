# gone

Gone is a wiki written in go. It's

* KISS,
* Convention over Configuration and
* designed with sysadmins in mind.

It might someday display plain text, markdown, godoc, manpages etc.
Files of unknown type are provided for download.
Plaintext files might be edited.

For permissions, the file system's access control is used as far as possible.

[![Build Status](https://travis-ci.org/fxnn/gone.svg?branch=master)](https://travis-ci.org/fxnn/gone)
[![GoDoc](https://godoc.org/github.com/fxnn/gone?status.svg)](https://godoc.org/github.com/fxnn/gone)


## Usage

Simply start the application.
The current working directory will now be served on port `8080`.
Append `?edit` to the URL to edit the content.
Append `?create` to the URL to create a non-existant file.


## Access Control

Users being not logged in have read and write permissions according to the
world access permissions of the file or containing directory.
This means, a file with `rwxrwx---` cannot be read or written by an
unauthenticated user.

### Currently unimplemented

Users can login.
The login information is configured in a good old `htpasswd` file.

By default, authenticated users can read and write all files that are readable
resp. writeable by the `gone` process.

Additionally, you can define user groups that must have read/write permissions
for anonymous/authenticated access.
Additional groups for authorized users might be supplied in a `htgroup` file.

For example, you can demand that access permissions for unauthenticated users
are only granted if the file has the `public` group.
You could also demand that authenticated users do not have additional rights,
unless they are granted special groups via `htgroup`.


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

