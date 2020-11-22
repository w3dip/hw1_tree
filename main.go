package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sort"
)

func filter(arr []os.FileInfo, cond func(os.FileInfo) bool) []os.FileInfo {
	result := []os.FileInfo{}
	for i := range arr {
	  if cond(arr[i]) {
		 result = append(result, arr[i])
	  }
	}
	return result
  }

func walk(out io.Writer, path string, printFiles bool, level int, parentIsLast bool, parentPrefix string) error {
	file, _ := os.Open(path)
	var items, _ = file.Readdir(0)
	if !printFiles {
		items = filter(items, func(val os.FileInfo) bool {
			return val.IsDir()
		})
	}

	items = filter(items, func(val os.FileInfo) bool {
		return !strings.HasPrefix(val.Name(), ".")
	})

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	if parentIsLast {
		parentPrefix += "\t"
	} else if level > 0 {
		parentPrefix += "│\t"
	}

	level++
	
	for indx, item := range items {
		name := item.Name()
		isDir := item.IsDir()
		isLastItem := indx == len(items) - 1
		
		var prefix string
		prefix += parentPrefix
		if isLastItem {
			prefix += "└───"
		} else {
			prefix += "├───"
		}

		if !isDir {
			if printFiles {
				var size string
				if item.Size() == 0 {
					size = "empty"
				} else {
					size = fmt.Sprintf("%vb", item.Size())
				}
				fmt.Fprintf(out, "%s%s (%s)\n", prefix, name, size)
			}
			continue
		}
		fmt.Fprintf(out, "%s%s\n", prefix, name)
		walk(out, filepath.Join(path, name), printFiles, level, isLastItem, parentPrefix)
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return walk(out, path, printFiles, 0, false, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
