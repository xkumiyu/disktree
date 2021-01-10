package disktree

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type Tree struct {
	Name     string
	Size     int64
	Files    int
	Dirs     int
	IsDir    bool
	Depth    int
	Children []Tree
}

func (t *Tree) Walk(path string) {
	ch := make(chan *Tree)
	files, _ := ioutil.ReadDir(path)
	for _, fi := range files {
		go func(file os.FileInfo) {
			st := Tree{
				Name:  file.Name(),
				IsDir: file.IsDir(),
				Depth: t.Depth + 1,
			}
			if st.IsDir {
				st.Dirs = 1
				st.Walk(filepath.Join(path, st.Name))
			} else {
				st.Files = 1
				st.Size = file.Size()
			}
			ch <- &st
		}(fi)
	}
	for range files {
		c := <-ch
		t.Size += c.Size
		t.Files += c.Files
		t.Dirs += c.Dirs
		t.Children = append(t.Children, *c)
	}
}

func (t *Tree) ReadableSize() string {
	size := float64(t.Size)
	i := 0
	for size >= 1000 && i < 5 {
		size /= 1000
		i++
	}
	var unit = [5]string{"B", "K", "M", "G", "T"}
	if size == 0 || size >= 10 {
		return fmt.Sprintf("%3.0f%s", size, unit[i])
	}
	return fmt.Sprintf("%1.1f%s", size, unit[i])
}

func (t *Tree) Sort(key string) {
	c := t.Children
	switch key {
	case "name":
		sort.SliceStable(c, func(i, j int) bool { return c[i].Name < c[j].Name })
	case "size":
		sort.SliceStable(c, func(i, j int) bool { return c[i].Size > c[j].Size })
	case "files":
		sort.SliceStable(c, func(i, j int) bool { return c[i].Files > c[j].Files })
	default:
		os.Exit(1)
	}
}

func (t *Tree) SizeFilter(minSize int64) {
	var newChildren []Tree
	for _, st := range t.Children {
		if st.Size > minSize {
			newChildren = append(newChildren, st)
		}
	}
	t.Children = newChildren
}
