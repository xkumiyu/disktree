package disktree_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/xkumiyu/disktree"
)

func TestDisktreeRun(t *testing.T) {
	tmp, err := ioutil.TempDir("", "dtreetest")
	if err != nil {
		t.Error("failed:", err)
	}
	defer os.RemoveAll(tmp)

	var buf bytes.Buffer
	d := disktree.New(tmp, -1, -1, "name", false, &buf)
	err = d.Run()
	if err != nil {
		t.Error("unexpected error:", err)
	}

	expected := fmt.Sprintf("  0B %s/\n\n0 directories, 0 files, 0 bytes\n", tmp)
	actual := buf.String()
	if expected != actual {
		t.Errorf("output wont %s but got %s", expected, actual)
	}
}
