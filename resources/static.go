package resources

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
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
	return nil, nil
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

	"/.viewer.html.swp": {
		local:   "resources/static/.viewer.html.swp",
		size:    12288,
		modtime: 1444227468,
		compressed: `
H4sIAAAJbogA/+yaz2oTQRjAp0WQIFVRRDwo03goluxO/pi0idtATQottLZq+o/Sw2R3kgxudpeZSZNQ
oqd68eAz+AIefANfwAfRox49+O0mxSpCbxbl+8GPb3b+fDPz5RR2m9mdtQ26YOcJcJ2Q5oNbO3c+nVxm
DiFuGHihIufiH0U6t1Cy8otW2coVC7YnPOlyIzy7E2ojeiqMhO2Rl+OErB+qFzrirmDtkGnlsrY0nV7T
dsMuaw2CALoDwZTQYU+5QjNtuJEuO5KiL5TdMV3//DMhCHKGnmlZizOkkM9l48f76Vl688b2RZ8KQRAE
QRAEQZC/iImmyCuI05Pne5M49VtEEARBEARBEARBEOTfhXuE7M4Q8gWM3/+f/v//eI2Q1+AJuA/ugU/B
BbAE3gWvgNPgh6uEvAffgW/BN2AE7oENcBWsgiWwCNrgbfAS+B32/QZ+nZzh88xFVgNBEARBEARBEOS/
wGHxp9PVuNEMvWGVpBzmySOIqePj5LPqWhgYEZjRCIZghLo+13op7Y6701XijBfGqQT3kgzaDH0R52Dz
9PCwSucZtCFBKtXlqi0DyxctU6FW1i6K7iPo7+QytJMHC+BDsAiW6PGvyyo0d7rAnuz/c0oLOqwW70p/
WKE17sumkhloBB5XPEOfi3YoMnQuiXR7bS5DNyMjuzC0rCT3M1TzQFtaKNmKN4gvlSSHKzizB7X6cmP5
ILmIk9yOmmEkltJGDAxztYY6wIiRBu4NhYu46YxGDht3EGdcGYhQUDro+gFUsGNMVGGs3+/b/YIdqjbL
lctlNojnxOn+MKHxbDycY/VGfdK0jIKDSyPDgPu2Z7z4J5mtb9Ya+1srNNlwa/vx+lqNpi3Gdgs1Fi+m
e6uNjXWoZ5Y2zqxnbOVJmvwIAAD//7L2wcwAMAAA
`,
	},

	"/editor.html": {
		local:   "resources/static/editor.html",
		size:    1137,
		modtime: 1442765421,
		compressed: `
H4sIAAAJbogA/4xTW2vbSBR+ln7FiWBfwkpjJSyLY9mQtQNbSJuQqLQl5GGsGVtDpBkxc2zZMfnvnYuj
NrSEvtg6l+9cvvNNcbK4mZffbq+gxraB28//XX+YQ5IS8uV8TsiiXMDX/8uP15BnIyg1lUagUJI2hFx9
SuKkRuwuCOn7PuvPM6XXpLwjO1crd+DjZ4o/ITOGLJnFceE77tpGmulv6uTj8TjAQzKnbBZHBQps+Oxw
yDqK9ctLQYIjtiGD+4YD7js+TZDvkFTGWGwUkVMoTh7mi8vy8gFOifWcwsH+RmnPl08C06XapUY8C7m+
gKXSjGvnmviUVj2/F38n1FHGvH80AWe3VK+FfDVf7MgWzfZhlNdgPuo82oezikvUPGS4lVLaiLXNcn6u
hzwXoppTOADYZTsVuAZhQMiaa4GchcWjhq9wmAhVN3xrsa59JGyFqFo3zXGcqBcMa+/46+1y2b+8nbgh
opWSmK5oK5r9BcyVNKqh5m9IrjeVYBTurQSgtNfp3Tw6+REJudx6WiuQSvl/ZTpa8clQ2JLMbf+zDsPa
/qyPjzO/V0H88WdxQYJQ4sJR6wRzkqawpEZU0NC92iCstGrhqDeDtHpSW65XjeqzSrWEkrP8n3E+Os8h
TR2eiS1UdhEr0sC6l1SxUrqFlmOt2DS5vbkvE6CVI32aDOL0mVEhZGfbBl3WgjEuE5C0tRZnAhPY0mYz
GCRghouGxMoyYJsnoFVvBzkbJe4JHL3+FRzzf+1oNsvWFQ6FDN3yoeO9N8ifYC4lu+O40fINGKhk8Or3
dQriiHG8EUucu0e4gzuMfcqz+HsAAAD//6KUEtlxBAAA
`,
	},

	"/viewer.html": {
		local:   "resources/static/viewer.html",
		size:    533,
		modtime: 1444227468,
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
