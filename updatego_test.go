package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMark(t *testing.T) {
	data := `
		// ErrorA represents "a"
		// TODO: add tests for 'ErrorA'
		ErrorA
		// ErrorB represents "b"
		ErrorB
		// ErrorC represents "c"
		ErrorD
		// ErrorE represents "d"
		`
	records := []Mark{{"ErrorA", "a", true}, {"ErrorB", "b", false}}

	rs := ParseMark(strings.Split(data, "\n"))
	assert.Equal(t, records, rs)
}
