package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/xkumiyu/disktree"
)

// Run ...
func Run() int {
	d, err := parseArgs()
	if err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	if err := d.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	return 0
}

func parseArgs() (*disktree.DiskTree, error) {
	var rootPath string

	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show the version and exit.")

	var noColor bool
	flag.BoolVar(&noColor, "no-color", false, "Disable colorization.")

	var maxDepth int
	flag.IntVar(&maxDepth, "max-depth", -1, "Show only to max-depth. -1 means infinity.")

	var sortKey string
	flag.StringVar(&sortKey, "sort", "name", "Select sort: name, size")

	// TODO: min-size
	// flag.IntVar(&minSize, "min-size", -1, "Show files/dirs larger than min-size.")

	flag.Parse()

	if showVersion {
		fmt.Printf("dtree version %s\n", disktree.Version)
		return nil, flag.ErrHelp
	}

	switch len(flag.Args()) {
	case 0:
		rootPath = "."
	case 1:
		rootPath = flag.Arg(0)
	default:
		return nil, errors.New("got unexpected extra argument")
	}

	if sortKey != "name" && sortKey != "size" {
		return nil, errors.New("invalid value for sort")
	}

	info, err := os.Stat(rootPath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, errors.New("not direcotry")
	}

	isColor := !noColor
	if !isTerminal() {
		isColor = false
	}

	return disktree.New(rootPath, maxDepth, sortKey, isColor, os.Stdout), nil
}

func isTerminal() bool {
	if info, _ := os.Stdout.Stat(); (info.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
