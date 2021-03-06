package main

import (
	"fmt"
	"strings"
)

func CutLeafs(leafs []Leaf) []Leaf {
	var ret []Leaf
	var last *Leaf
	last = nil
	for i := len(leafs); i > 0; i-- {
		leaf := leafs[i-1]
		if leaf.RFCRecord() {
			ret = append([]Leaf{leaf}, ret...)
		} else {
			if last != nil && leaf.Compare(*last) > 0 {
				ret = append([]Leaf{leaf}, ret...)
			}
		}
		last = &leaf
	}

	return ret
}

// var: poststartRef
// name: PostStart
// ref : config.md#poststart
type OutputRef struct {
	Var  string
	Name string
	Ref  string
}

func NewOutputRef(base string, leaf Leaf) OutputRef {
	var o OutputRef
	o.Var = ToCamel(leaf.Key, false)
	o.Name = leaf.Key
	o.Ref = fmt.Sprintf("%s#%s", base, ToRef(leaf.Key, false))
	return o
}

func (o *OutputRef) Debug(prefix string) {
	fmt.Printf("%s --- Ref --- %s %s %s\n", prefix, o.Var, o.Name, o.Ref)
}

// Keys:  Root, Path
// Value: "A directory MUST exist at the path declared by the field."
// RFC: Must
// Title: Root
type OutputRFC struct {
	Keys  []string
	Value string
	Ref   OutputRef
	RFC   string
	Title string
}

func (o *OutputRFC) Debug(prefix string) {
	fmt.Printf("%s --- RFC --- %s %s %s\n", prefix, strings.Join(o.Keys, "#"), o.Value, o.Ref)
}
func NewOutputRFC(base string, leafs []Leaf, cur int, title map[int]OutputRef) OutputRFC {
	var o OutputRFC
	curLevel := leafs[cur].Level
	o.Value = leafs[cur].Value
	o.RFC = ToCamel(strings.ToLower(leafs[cur].RFC), true)
	for i := cur; i >= 0; i-- {
		// Get Title, finish
		if leafs[i].Type == "title" {
			key := leafs[i].GetKey()
			o.Title = key
			o.Keys = append([]string{key}, o.Keys...)
			outputRef := NewOutputRef(base, leafs[i])
			if _, ok := title[i]; !ok {
				title[i] = outputRef
			}
			o.Ref = outputRef
			return o
		}

		// small means higher..
		if i == cur || leafs[i].Level < curLevel {
			curLevel = leafs[i].Level
			key := leafs[i].GetKey()
			o.Keys = append([]string{key}, o.Keys...)
		}
	}
	panic("Invalid loop?")
	return o
}

func OutputLeafs(base string, leafs []Leaf) ([]OutputRFC, []OutputRef) {
	var outputRFCs []OutputRFC
	var outputRefs []OutputRef
	usefulTitle := make(map[int]OutputRef)
	for i := len(leafs) - 1; i >= 0; i-- {
		leaf := leafs[i]
		if leaf.RFCRecord() {
			outputRFC := NewOutputRFC(base, leafs, i, usefulTitle)
			outputRFCs = append([]OutputRFC{outputRFC}, outputRFCs...)
		}
	}

	for i, l := range leafs {
		if l.Type == "title" {
			if _, ok := usefulTitle[i]; ok {
				outputRefs = append(outputRefs, usefulTitle[i])
			}
		}
	}
	return outputRFCs, outputRefs
}

func ToRef(value string, begin bool) string {
	name := strings.TrimSpace(value)

	ref := strings.Replace(name, "(", "", -1)
	ref = strings.Replace(ref, ")", "", -1)
	ref = strings.Replace(ref, " ", "-", -1)
	if !begin {
		ref = strings.ToLower(ref)
	}
	return ref
}

func ToCamel(value string, begin bool) string {
	value = ToRef(value, begin)
	var ret string
	last := begin
	for _, r := range value {
		switch r {
		case '-':
			last = true
		default:
			if last {
				ret = ret + strings.ToUpper(string(r))
				last = false
			} else {
				ret = ret + string(r)
			}
		}
	}
	return ret
}
