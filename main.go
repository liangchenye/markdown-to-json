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
	rfcs, refs := OutputLeafs(base, cleafs)

	ToGoTemplate(base, rfcs, refs)
}
