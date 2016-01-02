package filestore

import (
	"io"
	"io/ioutil"
)

type ioUtil struct {
	*errStore
}

func newIOUtil(errStore *errStore) *ioUtil {
	return &ioUtil{errStore}
}

func (i *ioUtil) readAllAndClose(readCloser io.ReadCloser) (result string) {
	if i.hasErr() {
		return ""
	}
	var buf []byte
	buf, err := ioutil.ReadAll(readCloser)
	i.setErr(err)
	readCloser.Close()
	return string(buf)
}

func (i *ioUtil) writeAllAndClose(writeCloser io.WriteCloser, content string) {
	if i.hasErr() {
		return
	}
	_, err := io.WriteString(writeCloser, content)
	i.setErr(err)
	writeCloser.Close()
}
