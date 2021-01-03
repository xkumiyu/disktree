package disktree_test

import (
	"bytes"
	"testing"

	"github.com/xkumiyu/disktree"
)

func TestDisktreeRun(t *testing.T) {
	var buf bytes.Buffer
	d := disktree.New(".", -1, "name", true, &buf)
	err := d.Run()
	if err != nil {
		t.Error("unexpected error:", err)
	}
}
