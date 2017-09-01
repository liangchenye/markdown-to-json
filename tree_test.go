package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCamel(t *testing.T) {
	cases := []struct {
		value    string
		expected string
	}{
		{"abcd", "abcd"},
		{"ab-cd", "abCd"},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, ToCamel(c.value))
	}
}
