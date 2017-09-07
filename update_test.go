package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateGoFile(t *testing.T) {
	expectedFile := "testdata/update.output"
	specErrFile := "testdata/error.go"
	specFile := "testdata/config.md"

	base := filepath.Base(specFile)
	data, _ := ioutil.ReadFile(specFile)

	markedData, _ := ioutil.ReadFile(specErrFile)

	leafs := GetLeafs(data)
	cleafs := CutLeafs(leafs)
	rfcs, refs := OutputLeafs(base, cleafs)
	unmarkedRfcs, unmarkedRefs := GetUnmarked(markedData, rfcs, refs)

	strs := UpdateGoFile(strings.Split(string(markedData), "\n"), unmarkedRfcs, unmarkedRefs)

	output := strings.Join(strs, "\n")
	expectedData, _ := ioutil.ReadFile(expectedFile)
	assert.Equal(t, expectedData, []byte(output))
}
