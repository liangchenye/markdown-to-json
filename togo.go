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
	errCodeTemplate = `
	// %s represents "%s"
        %s`
	errCodeEnd = `
)
`
)

var (
	refBegin        = `var (`
	refItemTemplate = `
	%sRef = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "%s"), nil
	}`
	refEnd = `
)
`
)

var (
	errRefBegin = `
var ociErrors = map[Code]errorTemplate{
	// %s`
	errRefTitleComment = `
	// %s`
	errRefTemplate = `
        %s: {Level: rfc2119.%s, Reference: %sRef},`

	errRefEnd = `
)
`
)

func ToGoTemplate(base string, rfcs []OutputRFC, refs []OutputRef) {
	// Error code
	fmt.Printf(errCodeBegin)
	for _, r := range rfcs {
		val := ToCamel(strings.Join(r.Keys, "-"), true)
		fmt.Printf(errCodeTemplate, val, strings.TrimSpace(r.Value), val)
	}
	fmt.Printf(errCodeEnd)

	fmt.Println("")

	// Refs
	fmt.Printf(refBegin)
	for _, r := range refs {
		fmt.Printf(refItemTemplate, r.Var, r.Ref)
	}
	fmt.Printf(refEnd)

	// Error ref
	fmt.Println("")
	fmt.Printf(errRefBegin, base)
	lastTitle := ""
	for _, r := range rfcs {
		if lastTitle != r.Title {
			fmt.Printf(errRefTitleComment, r.Title)
			lastTitle = r.Title
		}
		val := ToCamel(strings.Join(r.Keys, "-"), true)
		fmt.Printf(errRefTemplate, val, r.RFC, r.Ref.Var)
	}
	fmt.Printf(errRefEnd)
}
