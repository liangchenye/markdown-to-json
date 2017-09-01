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
	fmt.Println(base)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := GetLines(data)
	clines := CutLines(lines)
	//	clines := lines
	for _, l := range clines {
		l.Debug("")
	}
}
