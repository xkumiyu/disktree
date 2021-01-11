package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"

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
	var (
		rootPath    string
		showVersion bool
		color       string
		maxDepth    int
		minSize     string
		sortKey     string
		procs       int
	)

	flag.BoolVar(&showVersion, "version", false, "Show the version and exit.")
	flag.BoolVar(&showVersion, "V", false, "Alias of version.")

	flag.IntVar(&maxDepth, "max-depth", -1, "Show only to max-depth. -1 means infinity.")
	flag.IntVar(&maxDepth, "L", -1, "Alias of max-depth.")

	flag.StringVar(&sortKey, "sort", "name", "Select sort: name, size, and files")
	flag.StringVar(&sortKey, "s", "name", "Alias of sort.")

	flag.StringVar(&minSize, "min-size", "-1", "Show only files/dirs larger than min-size.")

	flag.StringVar(&color, "color", "auto", "Use terminal colors: auto, always, and never")

	flag.IntVar(&procs, "procs", -1, "The number of max processes. -1 means the number of logical CPUs.")

	flag.Parse()

	if showVersion {
		fmt.Printf("dtree version %s\n", disktree.Version)
		return nil, flag.ErrHelp
	}

	rootPath, err := parseRootPath()
	if err != nil {
		return nil, err
	}

	if sortKey != "name" && sortKey != "size" && sortKey != "files" {
		return nil, errors.New("invalid value for sort")
	}

	intMinSize, err := parseMinSize(minSize)
	if err != nil {
		return nil, err
	}

	isColor, err := parseColor(color)
	if err != nil {
		return nil, err
	}

	if procs != 1 {
		runtime.GOMAXPROCS(procs)
	}

	return disktree.New(rootPath, maxDepth, intMinSize, sortKey, isColor, os.Stdout), nil
}

func parseRootPath() (string, error) {
	var rootPath string
	switch len(flag.Args()) {
	case 0:
		rootPath = "."
	case 1:
		rootPath = flag.Arg(0)
	default:
		return "", errors.New("got unexpected extra argument")
	}

	info, err := os.Stat(rootPath)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		return "", errors.New("not direcotry")
	}

	return rootPath, nil
}

func parseColor(color string) (bool, error) {
	var isColor bool
	switch color {
	case "auto", "automatic":
		isColor = isTerminal()
	case "always", "on", "yes":
		isColor = true
	case "never", "off", "no":
		isColor = false
	default:
		return false, errors.New("invalid value for color")
	}
	return isColor, nil
}

func parseMinSize(minSize string) (int64, error) {
	var intMinSize int64
	r := regexp.MustCompile(`^(-?\d+)(B|K|M|G|T)?$`)
	m := r.FindAllStringSubmatch(minSize, -1)
	if len(m) != 1 {
		return -1, errors.New("invalid value for min-size")
	}
	intMinSize, _ = strconv.ParseInt(m[0][1], 10, 64)
	switch m[0][2] {
	case "K":
		intMinSize *= 1000
	case "M":
		intMinSize *= 1000 * 1000
	case "G":
		intMinSize *= 1000 * 1000 * 1000
	case "T":
		intMinSize *= 1000 * 1000 * 1000 * 1000
	}
	return intMinSize, nil
}

func isTerminal() bool {
	if info, _ := os.Stdout.Stat(); (info.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}
