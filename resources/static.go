package resources

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDir struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	data []byte
	once sync.Once
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDir) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	var fileInfos = make([]os.FileInfo, 0)
	for k, v := range _escData {
		if isFileInDir(k, f.name) {
			fileInfos = append(fileInfos, v)
		}
	}
	return fileInfos, nil
}

func isFileInDir(f string, d string) bool {
	if f == d {
		return false
	}

	// is f in d?
	if !strings.HasPrefix(f, d) {
		return false
	}

	// is f not in a subdir of d?
	var lastAllowedIndex = strings.LastIndexAny(d, "/")
	return strings.LastIndexAny(f, "/") == lastAllowedIndex
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDir{fs: _escLocal, name: name}
	}
	return _escDir{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		return ioutil.ReadAll(f)
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/editor.html": {
		local:   "resources/static/editor.html",
		size:    1206,
		modtime: 1451490921,
		compressed: `
H4sIAAAJbogA/4xUW2vjRhR+ln7FyUBfQqWxE0pxLLukdqCFtAmJSltCHsaakTVEmhGjY8uO8X/fuTja
hF3Cvlg6l+9cvvPJ2dnybpH/f38DFTY13P/z++2fCyAJpf9eLihd5kv474/8r1sYpyPIDVOdRKkVqym9
+ZvEpEJsryjt+z7tL1Nt1jR/oDtXa+zAp9cE3yFTjpzM4zjzHXdNrbrZd+qMJ5NJgIdkwfg8jjKUWIv5
4ZC2DKvjMaPBEdtQh/taAO5bMSModkiLrrPYKKLnkJ09LZbX+fUTnFPrOYeD/Y2SXqxeJCYrvUs6+SrV
+gpW2nBhnGvqUxr9+ln8k1DLOPf+0RSc3TCzlurNPNqRLZrvwyhvwfGo9WgfTguh0IiQ4VZKWC3XNsv5
hRnyXIgZweAAYJdtdeAaZAdSVcJIFDwsHtWixGEi1O3wbuS68pGwFaJu3DSncaJecqy846ePy6W/imbq
hohKrTApWSPr/RUstOp0zbqfgdxuCskZPFoJQG6v07t5DPkaCbnCehorkEL7p+5aVojpUNiSLGz/ixbD
2v6sz89zv1dG/fHncUaDUOLMUesEc5YksGKdLKBme71BKI1u4KS3DlnxorfClLXu00I3lNGL8S+T8ehy
DEni8FxuobCLWJEG1r2kslKbBhqBleYzcn/3mBNghSN9RgZx+swok6q1bYMuK8m5UAQUa6wluEQCW1Zv
BoMGzHDRkFhYBmxzAkb3dpCLEXGfwMnrv4JT/rcdu82qcYVDoY5txdDx0Rv0RzDXij8I3Bj1AQxMcXjz
hzqHgywhdbsc3Y1sYQaVEeU7Wn7johYoyHzpnxllJ6RQ3IMy6th15FPLvjtqOKa7rv0/mMdfAgAA///2
kUldtgQAAA==
`,
	},

	"/viewer.html": {
		local:   "resources/static/viewer.html",
		size:    533,
		modtime: 1449873698,
		compressed: `
H4sIAAAJbogA/2xQXWvbQBB89v2KjV4CQdLGdVOwexG4cqCBtAmtQltCHi7W2To4faBbKguj/949yQU/
5GHQsjczqxl5sXlMsz9Pd1BQaeHp+cvDfQpBhPhrkSJusg38/pp9e4B5fA1ZqypnyNSVsoh33wMRFETN
CrHrurhbxHW7x+wHHrzX3ItPY0RnyjinPEiEkOPFQ2krd/uOz3y5XE7yiaxVnoiZJENWJ8dj3CgqhkHi
tBD85Ki3Gqhv9G1A+kC4dY61sxlegbx4STfrbP0CV8ibtzrv4cjDbFdXFO1UaWy/glRZ89aakIcqV60K
4afe1zqEy/ELz/eXITw2ZEp+WrdG2RAcJ4ucbs3uM/sNjHjLnrqi6UCp2r2pVlzgjS7/U4p5CMUHxoLx
kXHD+HQuiKze0Qqi63MZB3l9TcYIEse4iZA4VSOkD+Urys1f2FrluNbTn4wtcGe+zXRaDYP3YKp3mJTe
igmJ+BcAAP//x9XRxhUCAAA=
`,
	},

	"/": {
		isDir: true,
		local: "resources/static",
	},
}
