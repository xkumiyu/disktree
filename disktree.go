package disktree

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

type node struct {
	name     string
	size     int64
	files    int
	dirs     int
	isdir    bool
	children []node
}

// Run ...
func (d *DiskTree) Run() error {
	n := node{
		name:     d.rootPath,
		children: d.walk(d.rootPath),
		isdir:    true,
	}
	for _, c := range n.children {
		n.size += c.size
		n.files += c.files
		n.dirs += c.dirs
	}
	d.print([]node{n}, "", 0)

	fmt.Fprintf(
		d.outWriter,
		"\n%d directories, %d files, %d bytes\n",
		n.dirs, n.files, n.size,
	)

	return nil
}

func (d *DiskTree) walk(path string) []node {
	var nodes []node

	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		n := node{
			name:  file.Name(),
			isdir: file.IsDir(),
		}

		if n.isdir {
			n.dirs = 1
			n.children = d.walk(filepath.Join(path, n.name))
			for _, c := range n.children {
				n.size += c.size
				n.files += c.files
				n.dirs += c.dirs
			}
		} else {
			n.files = 1
			n.size = file.Size()
		}

		nodes = append(nodes, n)
	}

	switch d.sortKey {
	case "size":
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].size > nodes[j].size
		})
	case "files":
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].files > nodes[j].files
		})
	case "name":
	default:
		os.Exit(1)
	}

	return nodes
}

func (d *DiskTree) print(nodes []node, basePrefix string, depth int) {
	if d.maxDepth != -1 && depth > d.maxDepth {
		return
	}

	if d.minSize != 1 {
		var filteredNodes []node
		for _, n := range nodes {
			if n.size > d.minSize {
				filteredNodes = append(filteredNodes, n)
			}
		}
		nodes = filteredNodes
	}

	for i, n := range nodes {
		size := readableSize(float64(n.size))
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
			if i == len(nodes)-1 {
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
			if i == len(nodes)-1 {
				nextPrefix += "    "
			} else {
				nextPrefix += "|   "
			}
		}
		d.print(n.children, nextPrefix, depth+1)
	}
}

func readableSize(size float64) string {
	i := 0
	for size > 1000 && i < 5 {
		size /= 1000
		i++
	}
	var unit = [5]string{"B", "K", "M", "G", "T"}
	if size == 0 || size >= 10 {
		return fmt.Sprintf("%3.0f%s", size, unit[i])
	}
	return fmt.Sprintf("%1.1f%s", size, unit[i])
}

func addColor(str string, color string) string {
	switch color {
	case "red":
		str = "\x1b[31m" + str + "\x1b[0m"
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
