package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var specFile string
	var specErrFile string
	var base string
	if len(os.Args) > 2 {
		specFile = os.Args[1]
		specErrFile = os.Args[2]
	} else {
		specFile = "testdata/config.md"
		specErrFile = "testdata/error.go"
	}

	base = filepath.Base(specFile)
	data, err := ioutil.ReadFile(specFile)
	if err != nil {
		fmt.Println(err, base, data)
		return
	}

	markedData, err := ioutil.ReadFile(specErrFile)
	if err != nil {
		fmt.Println(err, markedData)
		return
	}

	leafs := GetLeafs(data)
	cleafs := CutLeafs(leafs)
	rfcs, refs := OutputLeafs(base, cleafs)
	unmarkedRfcs, unmarkedRefs := GetUnmarked(markedData, rfcs, refs)

	strs := UpdateGoFile(strings.Split(string(markedData), "\n"), unmarkedRfcs, unmarkedRefs)

	output := strings.Join(strs, "\n")

	fmt.Println(output)
}
