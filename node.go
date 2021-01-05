package disktree

import (
	"fmt"
	"os"
	"sort"
)

type Nodes []Node

type Node struct {
	name     string
	size     int64
	files    int
	dirs     int
	isdir    bool
	children Nodes
}

func (n *Node) totalize() {
	for _, c := range n.children {
		n.size += c.size
		n.files += c.files
		n.dirs += c.dirs
	}
}

func (n *Node) readableSize() string {
	size := float64(n.size)
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

func (ns Nodes) Sort(key string) {
	switch key {
	case "size":
		sort.SliceStable(ns, func(i, j int) bool { return ns[i].size > ns[j].size })
	case "files":
		sort.SliceStable(ns, func(i, j int) bool { return ns[i].files > ns[j].files })
	default:
		os.Exit(1)
	}
}

func (ns Nodes) SizeFilter(minSize int64) Nodes {
	var newNodes Nodes
	for _, n := range ns {
		if n.size > minSize {
			newNodes = append(newNodes, n)
		}
	}
	return newNodes
}
