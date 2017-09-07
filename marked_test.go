package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUnmarked(t *testing.T) {
	data := `
		// ErrorA represents "a"
		// TODO: add tests for 'ErrorA'
		ErrorA
		// ErrorB represents "b"
		ErrorB
		// ErrorC represents "c"
		ErrorD
		// ErrorE represents "d"

		aRef = func(version string) (reference string, err error) {
		        return fmt.Sprintf(referenceTemplate, version, "config.md#specification-version"), nil
		}
		`
	markedRfcs := []OutputRFC{{Value: "a"}, {Value: "x"}}
	markedRefs := []OutputRef{{Var: "a"}, {Var: "x"}}

	expectedRfcs := []OutputRFC{{Value: "x"}}
	expectedRefs := []OutputRef{{Var: "xRef"}}
	retRfcs, retRefs := GetUnmarked([]byte(data), markedRfcs, markedRefs)

	assert.Equal(t, expectedRfcs, retRfcs)
	assert.Equal(t, expectedRefs, retRefs)
}

func TestParseMarked(t *testing.T) {
	data := `
		// ErrorA represents "a"
		// TODO: add tests for 'ErrorA'
		ErrorA
		// ErrorB represents "b"
		ErrorB
		// ErrorC represents "c"
		ErrorD
		// ErrorE represents "d"

		aRef = func(version string) (reference string, err error) {
		        return fmt.Sprintf(referenceTemplate, version, "config.md#specification-version"), nil
		}
		`
	rfcs := []MarkedRFC{{"ErrorA", "a", true}, {"ErrorB", "b", false}}
	refs := []string{"aRef"}

	retRfcs, retRefs := ParseMarked(strings.Split(data, "\n"))
	assert.Equal(t, rfcs, retRfcs)
	assert.Equal(t, refs, retRefs)
}
