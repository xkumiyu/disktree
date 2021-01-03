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
const Version = "0.1.1dev"

// DiskTree ...
type DiskTree struct {
	rootPath  string
	maxDepth  int
	sortKey   string
	isColor   bool
	outWriter io.Writer
}

// New DiskTree
func New(
	rootPath string,
	maxDepth int,
	sortKey string,
	isColor bool,
	outWriter io.Writer,
) *DiskTree {
	d := &DiskTree{
		rootPath:  rootPath,
		maxDepth:  maxDepth,
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
	}
	d.print([]node{n}, "", 0)

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
			n.children = d.walk(filepath.Join(path, n.name))
			for _, c := range n.children {
				n.size += c.size
				n.files += c.files
			}
		} else {
			n.files = 1
			n.size = file.Size()
		}

		nodes = append(nodes, n)
	}

	if d.sortKey == "size" {
		sort.SliceStable(nodes, func(i, j int) bool { return nodes[i].size > nodes[j].size })
	}

	return nodes
}

func (d *DiskTree) print(nodes []node, basePrefix string, depth int) {
	if d.maxDepth != -1 && depth > d.maxDepth {
		return
	}

	// TODO: filter files using min-size

	for i, node := range nodes {
		size := readableSize(float64(node.size))
		if d.isColor {
			size = addColor(size, "green")
		}
		name := node.name
		if node.isdir {
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
		if node.isdir && node.files > 0 {
			num := fmt.Sprintf("[%d files]", node.files)
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
		d.print(node.children, nextPrefix, depth+1)
	}
}

func readableSize(size float64) string {
	i := 0
	for size > 1000 && i < 5 {
		size /= 1024
		i++
	}
	var unit = [5]string{"B", "K", "M", "G", "T"}
	if size < 10 {
		return fmt.Sprintf("%1.1f%s", size, unit[i])
	}
	return fmt.Sprintf("%3.0f%s", size, unit[i])
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
