package main

import (
	"fmt"
	"strings"
)

var (
	headTemplate = `package specerror

import (
    "fmt"

    rfc2119 "github.com/opencontainers/runtime-tools/error"
)
`

	codeBeginTemplate = `
// define error codes
const (
`
	codeItemTemplate = `	// %s represents "%s"
        %s = "%s"
`
	codeEndTemplate = `)
`

	refBeginTemplate = `var (
`
	refItemTemplate = `	%s = func(version string) (reference string, err error) {
		return fmt.Sprintf(referenceTemplate, version, "%s"), nil
	}
`
	refEndTemplate = `)`

	regBeginTemplate = `
func init() {`
	regItemTemplate = `
	registOCIError(%s, rfc2119.%s, %sRef)`
	regEndTemplate = `
}
`
)

func GetOutputCodeContent(rfcs []OutputRFC, remind string) []string {
	var ret []string
	for i, r := range rfcs {
		val := fmt.Sprintf("%s%d", ToCamel(strings.Join(r.Keys, "-"), true), i)
		if remind != "" {
			ret = append(ret, remind)
		}
		ret = append(ret, fmt.Sprintf(codeItemTemplate, val, strings.TrimSpace(r.Value), val, strings.TrimSpace(r.Value)))
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

func GetOutputRegContent(rfcs []OutputRFC, remind string) []string {
	var ret []string
	for i, r := range rfcs {
		if remind != "" {
			ret = append(ret, remind)
		}
		val := fmt.Sprintf("%s%d", ToCamel(strings.Join(r.Keys, "-"), true), i)
		ret = append(ret, fmt.Sprintf(regItemTemplate, val, r.RFC, r.Ref.Var))
	}
	return ret
}

func ToGoTemplate(base string, rfcs []OutputRFC, refs []OutputRef) {
	fmt.Printf(base)

	// Error code
	fmt.Printf(codeBeginTemplate)
	fmt.Printf(strings.Join(GetOutputCodeContent(rfcs, ""), ""))
	fmt.Printf(codeEndTemplate)

	fmt.Println("")

	// Refs
	fmt.Printf(refBeginTemplate)
	fmt.Printf(strings.Join(GetOutputRefContent(refs, ""), ""))
	fmt.Printf(refEndTemplate)

	// Error map
	fmt.Println("")
	fmt.Printf(regBeginTemplate)
	fmt.Printf(strings.Join(GetOutputRegContent(rfcs, ""), ""))
	fmt.Printf(regEndTemplate)
}
