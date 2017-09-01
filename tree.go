package main

import (
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
			} else {
				//				line.Debug("drop")
				//				if last != nil {
				//					last.Debug("\t\t")
				//				}
			}
		}
		last = &line
	}

	return ret
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
