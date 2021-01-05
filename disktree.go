package disktree

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Version of DiskTree
const Version = "0.2.0-beta"

// DiskTree ...
type DiskTree struct {
	rootPath  string
	maxDepth  int
	minSize   int64
	sortKey   string
	isColor   bool
	outWriter io.Writer
}

// New DiskTree
func New(
	rootPath string,
	maxDepth int,
	minSize int64,
	sortKey string,
	isColor bool,
	outWriter io.Writer,
) *DiskTree {
	d := &DiskTree{
		rootPath:  rootPath,
		maxDepth:  maxDepth,
		minSize:   minSize,
		sortKey:   sortKey,
		isColor:   isColor,
		outWriter: outWriter,
	}
	return d
}

// Run ...
func (d *DiskTree) Run() error {
	n := Node{
		name:     d.rootPath,
		children: d.walk(d.rootPath),
		isdir:    true,
	}
	n.totalize()
	d.print([]Node{n}, "", 0)

	fmt.Fprintf(
		d.outWriter,
		"\n%d directories, %d files, %d bytes\n",
		n.dirs, n.files, n.size,
	)

	return nil
}

func (d *DiskTree) walk(path string) Nodes {
	var ns Nodes

	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		n := Node{
			name:  file.Name(),
			isdir: file.IsDir(),
		}
		if n.isdir {
			n.dirs = 1
			n.children = d.walk(filepath.Join(path, n.name))
			n.totalize()
		} else {
			n.files = 1
			n.size = file.Size()
		}
		ns = append(ns, n)
	}
	return ns
}

func (d *DiskTree) print(ns Nodes, basePrefix string, depth int) {
	if d.maxDepth != -1 && depth > d.maxDepth {
		return
	}
	if d.minSize != 1 {
		ns = ns.SizeFilter(d.minSize)
	}
	if d.sortKey != "name" {
		ns.Sort(d.sortKey)
	}

	for i, n := range ns {
		size := n.readableSize()
		if d.isColor {
			size = addColor(size, "green")
		}
		name := n.name
		if n.isdir {
			if d.isColor {
				name = addColor(name, "blue")
			}
			name += "/"
		}
		body := fmt.Sprintf("%s %s", size, name)

		prefix := basePrefix
		if depth > 0 {
			if i == len(ns)-1 {
				prefix += "`-- "
			} else {
				prefix += "|-- "
			}
		}

		suffix := ""
		if n.isdir && n.files > 0 {
			num := fmt.Sprintf("[%d files]", n.files)
			if d.isColor {
				num = addColor(num, "yellow")
			}
			suffix = " " + num
		}

		fmt.Fprintln(d.outWriter, prefix+body+suffix)

		nextPrefix := basePrefix
		if depth > 0 {
			if i == len(ns)-1 {
				nextPrefix += "    "
			} else {
				nextPrefix += "|   "
			}
		}
		d.print(n.children, nextPrefix, depth+1)
	}
}

func addColor(str string, color string) string {
	switch color {
	case "green":
		str = "\x1b[32m" + str + "\x1b[0m"
	case "yellow":
		str = "\x1b[33m" + str + "\x1b[0m"
	case "blue":
		str = "\x1b[34m" + str + "\x1b[0m"
	default:
		os.Exit(1)
	}
	return str
}
