package main

import (
	"fmt"
	"strings"
)

var (
	errCodeBegin = `
// Code represents the spec violation, enumerating both
// configuration violations and runtime violations.
type Code int

const (
        // NonError represents that an input is not an error
        NonError Code = iota
        // NonRFCError represents that an error is not a rfc2119 error
        NonRFCError

`
	errCodeTemplate = `	// %s represents "%s"
        %s
`
	errCodeEnd = `)
`
)

var (
	refBegin = `var (
`
	refItemTemplate = `	%sRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "%s"), nil
	}
`
	refEnd = `)
`
)

var (
	errMapBegin = `var ociErrors = map[Code]errorTemplate{
	// %s
`
	errMapTitleComment = `	// %s
`
	errMapTemplate = `	%s: {Level: rfc2119.%s, Reference: %sRef},
`
	errMapEnd = `)
`
)

func GetOutputCodeContent(rfcs []OutputRFC, remind string) []string {
	var ret []string
	for _, r := range rfcs {
		val := ToCamel(strings.Join(r.Keys, "-"), true)
		if remind != "" {
			ret = append(ret, remind)
		}
		ret = append(ret, fmt.Sprintf(errCodeTemplate, val, strings.TrimSpace(r.Value), val))
	}
	return ret
}

func GetOutputRefContent(refs []OutputRef, remind string) []string {
	var ret []string
	for _, r := range refs {
		if remind != "" {
			ret = append(ret, remind)
		}
		ret = append(ret, fmt.Sprintf(refItemTemplate, r.Var, r.Ref))
	}
	return ret
}

func GetOutputMapContent(rfcs []OutputRFC, remind string) []string {
	var ret []string
	lastTitle := ""
	for _, r := range rfcs {
		if lastTitle != r.Title {
			ret = append(ret, fmt.Sprintf(errMapTitleComment, r.Title))
			lastTitle = r.Title
		}
		if remind != "" {
			ret = append(ret, remind)
		}
		val := ToCamel(strings.Join(r.Keys, "-"), true)
		ret = append(ret, fmt.Sprintf(errMapTemplate, val, r.RFC, r.Ref.Var))
	}
	return ret
}

func ToGoTemplate(base string, rfcs []OutputRFC, refs []OutputRef) {
	// Error code
	fmt.Printf(errCodeBegin)
	fmt.Printf(strings.Join(GetOutputCodeContent(rfcs, ""), ""))
	fmt.Printf(errCodeEnd)

	fmt.Println("")

	// Refs
	fmt.Printf(refBegin)
	fmt.Printf(strings.Join(GetOutputRefContent(refs, ""), ""))
	fmt.Printf(refEnd)

	// Error map
	fmt.Println("")
	fmt.Printf(errMapBegin, base)
	fmt.Printf(strings.Join(GetOutputMapContent(rfcs, ""), ""))
	fmt.Printf(errMapEnd)
}
