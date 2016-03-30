package templates

import (
	"os"
	"testing"

	"github.com/fxnn/gopath"
)

func TestExportTemplates(t *testing.T) {
	var wd, _ = os.Getwd()
	var sut = NewStaticLoader()
	var targetPath = gopath.FromPath(wd).JoinPath("gone_test")
	defer os.RemoveAll(targetPath.Path())
	var err = sut.WriteAllTemplates(targetPath)
	if err != nil {
		t.Fatalf("Couldn't write templates: %v", err)
	}
}
