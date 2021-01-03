package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type obj struct {
	name     string
	size     int64
	children []obj
	num      int
	isdir    bool
}

func main() {
	var root string

	flag.Parse()
	switch len(flag.Args()) {
	case 0:
		root = "."
	case 1:
		root = flag.Arg(0)
	default:
		os.Exit(1)
	}

	o := obj{
		name:     root,
		size:     0,
		children: walk(root),
		num:      0,
		isdir:    true,
	}
	for _, c := range o.children {
		o.size += c.size
		o.num += c.num
	}

	print([]obj{o}, "", 0)
}

func print(files []obj, prefix string, depth int) {
	for i, file := range files {
		size := addColor(readableSize(file.size), "green")
		name := file.name
		if file.isdir {
			name = addColor(name, "blue") + "/"
		}
		body := fmt.Sprintf("%s %s", size, name)
		if depth > 0 {
			if i == len(files)-1 {
				body = "`-- " + body
			} else {
				body = "|-- " + body
			}
		}

		suffix := ""
		if file.isdir && file.num > 0 {
			num := addColor(fmt.Sprintf("[%d files]", file.num), "yellow")
			suffix = " " + num
		}

		fmt.Printf("%s%s%s\n", prefix, body, suffix)

		if depth > 0 {
			if i == len(files)-1 {
				print(file.children, prefix+"    ", depth+1)
			} else {
				print(file.children, prefix+"|   ", depth+1)
			}
		} else {
			print(file.children, prefix, depth+1)
		}
	}
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

func walk(path string) []obj {
	var arr []obj

	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		o := obj{
			name:  file.Name(),
			size:  file.Size(),
			num:   1,
			isdir: file.IsDir(),
		}
		if file.IsDir() {
			o.num = 0
			o.children = walk(filepath.Join(path, o.name))
			o.size = 0
			for _, c := range o.children {
				o.size += c.size
				o.num += c.num
			}
		}
		arr = append(arr, o)
	}

	// sort.SliceStable(arr, func(i, j int) bool { return arr[i].size > arr[j].size })

	return arr
}

func readableSize(size int64) string {
	if size > 1000*1000*1000*1000 {
		return fmt.Sprintf("%3dT", size/(1024*1024*1024*1024))
	} else if size > 1000*1000*1000 {
		return fmt.Sprintf("%3dG", size/(1024*1024*1024))
	} else if size > 1000*1000 {
		return fmt.Sprintf("%3dM", size/(1024*1024))
	} else if size > 1000 {
		return fmt.Sprintf("%3dK", size/1024)
	} else {
		return fmt.Sprintf("%3dB", size)
	}
}
