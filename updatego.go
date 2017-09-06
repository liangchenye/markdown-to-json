package main

import (
	"strings"
)

var (
	codeComment = match("// (\\w+) represents \"([\\S\\s]+)\"")
	codeTodo    = match("TODO: add tests for")
)

type Mark struct {
	ErrorCode string
	Represent string
	Todo      bool
}

func ParseMark(lines []string) []Mark {
	var records []Mark
	var r *Mark
	status := "none"
	for _, l := range lines {
		switch status {
		case "none":
			vs := codeComment.FindStringSubmatch(l)
			if len(vs) != 3 {
				continue
			}
			r = &Mark{}
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
				records = append(records, *r)
			}
		case "code":
			if strings.TrimSpace(l) == r.ErrorCode {
				status = "none"
				records = append(records, *r)
			}
		}
	}

	return records
}
