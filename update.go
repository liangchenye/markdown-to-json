package main

var (
	codeBeginReg = match("^// define error codes")
	codeEndReg   = match("^\\)")

	refBeginReg = match("Ref = func")
	refEndReg   = match("^\\)")

	errMapBeginReg = match("^\tregistOCIError")
	errMapEndReg   = match("^\\}")
)

var (
	remindTemplate = `	// TODO: newly added code, need to update.`
)

func UpdateGoFile(lines []string, rfcs []OutputRFC, refs []OutputRef) []string {
	var ret []string
	status := "none"
	for _, l := range lines {
		switch status {
		case "none":
			if codeBeginReg.MatchString(l) {
				status = "codeBegin"
			} else if refBeginReg.MatchString(l) {
				status = "refBegin"
			} else if errMapBeginReg.MatchString(l) {
				status = "mapBegin"
			}
		case "codeBegin":
			if codeEndReg.MatchString(l) {
				ret = append(ret, GetOutputCodeContent(rfcs, remindTemplate)...)
				status = "none"
			}
		case "refBegin":
			if refEndReg.MatchString(l) {
				ret = append(ret, GetOutputRefContent(refs, remindTemplate)...)
				status = "none"
			}
		case "mapBegin":
			if errMapEndReg.MatchString(l) {
				ret = append(ret, GetOutputRegContent(rfcs, remindTemplate)...)
				status = "none"
			}
		}
		ret = append(ret, l)
	}
	return ret
}
