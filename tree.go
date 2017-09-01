package main

import (
	"fmt"
	"strings"
)

func CutLines(lines []Line) []Line {
	var ret []Line
	var last *Line
	last = nil
	for i := len(lines); i > 0; i-- {
		line := lines[i-1]
		if line.RFCRecord() {
			ret = append([]Line{line}, ret...)
		} else {
			if last != nil && line.Compare(*last) > 0 {
				ret = append([]Line{line}, ret...)
			}
		}
		last = &line
	}

	return ret
}

type OutputRef struct {
	Var  string
	Name string
	Ref  string
}

func NewOutputRef(base string, line Line) OutputRef {
	var o OutputRef
	o.Var = ToCamel(line.T.Ref)
	o.Name = line.T.Name
	o.Ref = fmt.Sprintf("%s#%s", base, line.T.Ref)
	return o
}

func (o *OutputRef) Debug(prefix string) {
	fmt.Printf("%s %s %s %s\n", prefix, o.Var, o.Name, o.Ref)
}

type OutputRFC struct {
	Keys  []string
	Value string
	Ref   OutputRef
}

func (o *OutputRFC) Debug(prefix string) {
	fmt.Printf("%s %s %s %s\n", prefix, strings.Join(o.Keys, "#"), o.Value, o.Ref)
}
func NewOutputRFC(base string, lines []Line, cur int, title map[int]OutputRef) OutputRFC {
	var o OutputRFC
	o.Value = lines[cur].I.Value
	for i := cur; i >= 0; i-- {
		key, isTitle := lines[i].GetKey()
		if key != "" {
			o.Keys = append([]string{key}, o.Keys...)
		}
		// When we get title, we get ref and thus to return
		if isTitle {
			outputRef := NewOutputRef(base, lines[i])
			title[i] = outputRef
			o.Ref = outputRef
			return o
		}
	}
	panic("Invalid loop?")
	return o
}

func OutputLines(base string, lines []Line) ([]OutputRFC, []OutputRef) {
	var outputRFCs []OutputRFC
	var outputRefs []OutputRef
	usefulTitle := make(map[int]OutputRef)
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if line.RFCRecord() {
			outputRFC := NewOutputRFC(base, lines, i, usefulTitle)
			outputRFCs = append([]OutputRFC{outputRFC}, outputRFCs...)
		}
	}

	for _, ref := range usefulTitle {
		outputRefs = append(outputRefs, ref)
	}
	return outputRFCs, outputRefs
}

func ToCamel(value string) string {
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
