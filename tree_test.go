package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToCamel(t *testing.T) {
	cases := []struct {
		value    string
		begin    bool
		expected string
	}{
		{"abcd", false, "abcd"},
		{"ab-cd", true, "AbCd"},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, ToCamel(c.value, c.begin))
	}
}
