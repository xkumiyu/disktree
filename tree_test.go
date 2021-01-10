package disktree_test

import (
	"testing"

	"github.com/xkumiyu/disktree"
)

func TestTreeReadableSize(t *testing.T) {
	cases := map[string]struct {
		size     int64
		expected string
	}{
		"0B":   {size: 0, expected: "  0B"},
		"100B": {size: 100, expected: "100B"},
		"1K":   {size: 1000, expected: "1.0K"},
		"10M":  {size: 10000000, expected: " 10M"},
		"1.2G": {size: 1234567890, expected: "1.2G"},
	}

	for n, tt := range cases {
		tt := tt
		t.Run(n, func(t *testing.T) {
			tree := disktree.Tree{Size: tt.size}
			if actual := tree.ReadableSize(); tt.expected != actual {
				t.Errorf("readable size wont %s but got %s", tt.expected, actual)
			}
		})
	}
}
