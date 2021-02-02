package disktree

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// Version of DiskTree
const Version = "0.3.0"

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

const clearLine = "\r\033[K"

// Run ...
func (d *DiskTree) Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(1 * time.Second)
		spin(ctx)
	}()

	t := Tree{Name: d.rootPath, IsDir: true}
	t.Walk(d.rootPath)

	// fmt.Print("\033[?25h")
	fmt.Print(clearLine)
	cancel()

	d.print(&t, "", false)
	fmt.Fprintf(
		d.outWriter,
		"\n%d directories, %d files, %d bytes\n",
		t.Dirs, t.Files, t.Size,
	)

	return nil
}

func (d *DiskTree) print(t *Tree, basePrefix string, isLeaf bool) {
	if d.maxDepth != -1 && t.Depth > d.maxDepth {
		return
	}
	if d.minSize != 1 {
		t.SizeFilter(d.minSize)
	}
	// TODO: only with all option
	ignoreDot := false
	if ignoreDot {
		t.ignoreFilter()
	}
	t.Sort(d.sortKey)

	size := t.ReadableSize()
	if d.isColor {
		size = addColor(size, "green")
	}
	name := t.Name
	if t.IsDir {
		if d.isColor {
			name = addColor(name, "blue")
		}
		name += "/"
	}
	body := fmt.Sprintf("%s %s", size, name)

	prefix := basePrefix
	nextPrefix := basePrefix
	if t.Depth > 0 {
		if isLeaf {
			prefix += "`-- "
			nextPrefix += "    "
		} else {
			prefix += "|-- "
			nextPrefix += "|   "
		}
	}

	suffix := ""
	if t.IsDir && t.Files > 0 {
		num := fmt.Sprintf("[%d files]", t.Files)
		if d.isColor {
			num = addColor(num, "yellow")
		}
		suffix = " " + num
	}

	fmt.Fprintln(d.outWriter, prefix+body+suffix)

	for i, st := range t.Children {
		var nextIsLeaf bool
		if i == len(t.Children)-1 {
			nextIsLeaf = true
		} else {
			nextIsLeaf = false
		}
		d.print(&st, nextPrefix, nextIsLeaf)
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

func spin(ctx context.Context) {
	frames := []rune(`⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏`)
	delay := 100 * time.Millisecond

	// fmt.Print("\033[?25l")
	for {
		for i := 0; i < len(frames); i++ {
			select {
			case <-ctx.Done():
				return
			default:
				fmt.Printf("\r%s Exploring...", string(frames[i]))
				time.Sleep(delay)
			}

		}
	}
}
