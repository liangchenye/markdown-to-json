package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var file string
	var base string
	if len(os.Args) > 1 {
		file = os.Args[1]
	} else {
		file = "testdata/config.md"
	}

	base = filepath.Base(file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	leafs := GetLeafs(data)
	cleafs := CutLeafs(leafs)
	//	cleafs := leafs
	rfcs, refs := OutputLeafs(base, cleafs)
	fmt.Println("RFCS ----")
	for _, l := range rfcs {
		l.Debug("")
	}
	fmt.Println("REFS----")
	for _, l := range refs {
		l.Debug("")
	}
}
