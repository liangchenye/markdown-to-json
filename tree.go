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

type OutputRef struct {
	Var  string
	Name string
	Ref  string
}

func NewOutputRef(base string, leaf Leaf) OutputRef {
	var o OutputRef
	o.Var = ToCamel(leaf.Key)
	o.Name = leaf.Key
	o.Ref = fmt.Sprintf("%s#%s", base, ToRef(leaf.Key))
	return o
}

func (o *OutputRef) Debug(prefix string) {
	fmt.Printf("%s --- Ref --- %s %s %s\n", prefix, o.Var, o.Name, o.Ref)
}

type OutputRFC struct {
	Keys  []string
	Value string
	Ref   OutputRef
}

func (o *OutputRFC) Debug(prefix string) {
	fmt.Printf("%s --- RFC --- %s %s %s\n", prefix, strings.Join(o.Keys, "#"), o.Value, o.Ref)
}
func NewOutputRFC(base string, leafs []Leaf, cur int, title map[int]OutputRef) OutputRFC {
	var o OutputRFC
	curLevel := leafs[cur].Level
	o.Value = leafs[cur].Value
	for i := cur; i >= 0; i-- {
		// Get Title, finish
		if leafs[i].Type == "title" {
			key := leafs[i].GetKey()
			o.Keys = append([]string{key}, o.Keys...)
			outputRef := NewOutputRef(base, leafs[i])
			title[i] = outputRef
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

	for _, ref := range usefulTitle {
		outputRefs = append(outputRefs, ref)
	}
	return outputRFCs, outputRefs
}

func ToRef(value string) string {
	name := strings.TrimSpace(value)

	ref := strings.Replace(name, "(", "", -1)
	ref = strings.Replace(ref, ")", "", -1)
	ref = strings.Replace(ref, " ", "-", -1)
	ref = strings.ToLower(ref)
	return ref
}

func ToCamel(value string) string {
	value = ToRef(value)
	var ret string
	last := false
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
