package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReg(t *testing.T) {
	cases := []struct {
		value    string
		title    Title
		expected bool
	}{
		{"##<a> b </a> HEllo world", Title{"##<a> b </a> HEllo world", 2, "HEllo world", "hello-world"}, true},
		{"##<a> b </a> hello (world)", Title{"##<a> b </a> hello (world)", 2, "hello (world)", "hello-world"}, true},
		{"#<a name=\"ab\"> /> hello world", Title{"#<a name=\"ab\"> /> hello world", 1, "hello world", "hello-world"}, true},
		{"<a>b c </a> hello world", Title{"", 0, "", ""}, false},
	}

	for _, c := range cases {
		title, err := NewTitle(c.value)
		assert.Equal(t, c.title, title)
		assert.Equal(t, c.expected, err == nil)
	}
}

func TestItem(t *testing.T) {
	cases := []struct {
		value string
		item  Item
	}{
		{"    **`myVersion`** it MUST", Item{"    **`myVersion`** it MUST", 4, "myVersion", "MUST"}},
		{"hello", Item{"hello", 0, "", ""}},
	}

	for _, c := range cases {
		item := NewItem(c.value)
		assert.Equal(t, c.item, item)
	}
}
