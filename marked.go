package main

import (
	"fmt"
	"strings"
)

var (
	codeRef     = match("(\\w+) = func\\(version string\\) \\(reference string, err error\\)")
	codeComment = match("// (\\w+) represents \"([\\S\\s]+)\"")
	codeTodo    = match("TODO: add tests for")
)

type MarkedRFC struct {
	ErrorCode string
	Represent string
	Todo      bool
}

func GetUnmarked(markedData []byte, rfcs []OutputRFC, refs []OutputRef) ([]OutputRFC, []OutputRef) {
	var retRfcs []OutputRFC
	var retRefs []OutputRef
	markedRfcs, markedRefs := ParseMarked(strings.Split(string(markedData), "\n"))
	markMap := make(map[string]bool)
	for _, m := range markedRfcs {
		markMap[m.Represent] = true
	}
	for _, m := range markedRefs {
		markMap[m] = true
	}

	for _, rfs := range rfcs {
		if _, ok := markMap[strings.TrimSpace(rfs.Value)]; !ok {
			retRfcs = append(retRfcs, rfs)
		}
	}

	for _, ref := range refs {
		if _, ok := markMap[strings.TrimSpace(ref.Var)]; !ok {
			retRefs = append(retRefs, ref)
		}
	}
	return retRfcs, retRefs
}

func ParseMarked(lines []string) ([]MarkedRFC, []string) {
	var markedRfcs []MarkedRFC
	var markedRefs []string
	var r *MarkedRFC
	status := "none"
	for _, l := range lines {
		// Get Ref
		refs := codeRef.FindStringSubmatch(l)
		if len(refs) == 2 {
			markedRefs = append(markedRefs, strings.TrimSpace(refs[1]))
			continue
		}
		switch status {
		case "none":

			vs := codeComment.FindStringSubmatch(l)
			if len(vs) != 3 {
				continue
			}
			r = &MarkedRFC{}
			r.ErrorCode = vs[1]
			r.Represent = vs[2]
			status = "comment"
		case "comment":
			vs := codeTodo.FindStringSubmatch(l)
			if len(vs) != 0 {
				r.Todo = true
				status = "code"
				continue
			}
			if strings.TrimSpace(l) == r.ErrorCode {
				status = "none"
				markedRfcs = append(markedRfcs, *r)
			}
		case "code":
			if strings.TrimSpace(l) == r.ErrorCode {
				status = "none"
				markedRfcs = append(markedRfcs, *r)
			}
		}
	}

	return markedRfcs, markedRefs
}
