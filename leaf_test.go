package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTitle(t *testing.T) {
	cases := []struct {
		value    string
		title    Leaf
		expected bool
	}{
		{"##<a> b </a> HEllo world", Leaf{"title", "##<a> b </a> HEllo world", 2, "HEllo world", "hello-world"}, true},
		{"##<a> b </a> hello (world)", Leaf{"title", "##<a> b </a> hello (world)", 2, "hello (world)", "hello-world"}, true},
		{"#<a name=\"ab\"> /> hello world", Leaf{"title", "#<a name=\"ab\"> /> hello world", 1, "hello world", "hello-world"}, true},
	}

	for _, c := range cases {
		title := NewLeaf(c.value)
		assert.Equal(t, c.expected, title != nil)
		if title != nil {
			assert.Equal(t, c.title, *title)
		}
	}
}

func TestItem(t *testing.T) {
	cases := []struct {
		value string
		item  Leaf
	}{
		{"    **`myVersion`** it MUST", Leaf{"item", "    **`myVersion`** it MUST", 4, "myVersion", "MUST"}},
		{"hello", Leaf{"item", "hello", 0, "", ""}},
	}

	for _, c := range cases {
		item := NewLeaf(c.value)
		assert.Equal(t, c.item, *item)
	}
}
