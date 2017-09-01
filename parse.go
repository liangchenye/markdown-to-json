package main

import (
	"fmt"
	"regexp"
	"strings"
)

const TitleWeight = 100

var (
	match = regexp.MustCompile
	// Limit Title to less than 6 levels
	titleReg    = match("([#]{1,6})(.*)")
	anchorReg   = match("\\<a[\\S\\s]+?\\</a\\>|<a[\\S\\s]+?/\\>")
	backtickReg = match("```[\\S\\s]+?```")
	rfcReg      = match("(MUST|SHOULD|REQUIRED)")
	keyReg      = match("\\*\\*`(\\w+)`\\*\\*")
)

// ## <a> abcd </a> Hello word
// In this case
// Value is '## <a> abcd </a> Hello word'
// Level is 2
// Name is 'Hello word'
// Ref is  'hello-word'
type Title struct {
	Value string
	Level int
	Name  string
	Ref   string
}

func GetTitleNameAndRef(value string) (string, string) {
	name := strings.TrimSpace(value)

	ref := strings.Replace(name, "(", "", -1)
	ref = strings.Replace(ref, ")", "", -1)
	ref = strings.Replace(ref, " ", "-", -1)
	ref = strings.ToLower(ref)
	return name, ref
}

func NewTitle(value string) (Title, error) {
	var title Title
	// trim all the anchors
	src := anchorReg.ReplaceAllString(value, "")
	ts := titleReg.FindStringSubmatch(src)
	if len(ts) != 3 {
		return Title{}, fmt.Errorf("'%s' is not valid title", value)
	}

	title.Value = value
	title.Level = len(ts[1])
	title.Name, title.Ref = GetTitleNameAndRef(ts[2])
	return title, nil
}

type Item struct {
	Value string
	// Item is the blank space numbers at the beginning
	//   the line which has more spaces, is the child most likely
	Level int
	// Could be empty
	Key string
	// FIXME: we only assume one, could be empty
	RFC string
}

func NewItem(value string) Item {
	var item Item

	item.Value = value
	for _, v := range value {
		if v == ' ' {
			item.Level++
		} else {
			break
		}
	}

	item.Level += TitleWeight

	vs := rfcReg.FindStringSubmatch(value)
	if len(vs) == 2 {
		item.RFC = vs[1]
	}

	vs = keyReg.FindStringSubmatch(value)
	if len(vs) == 2 {
		item.Key = vs[1]
	}

	return item
}

func (i *Item) Useful() bool {
	if i.Key != "" || i.RFC != "" {
		return true
	}

	return false
}

type Line struct {
	T *Title
	I *Item
}

func NewLine(value string) *Line {
	var l Line
	l.T = nil
	l.I = nil
	if strings.HasPrefix(value, "#") {
		if title, err := NewTitle(value); err == nil {
			l.T = &title
			return &l
		}
	} else {
		item := NewItem(value)
		if item.Useful() {
			l.I = &item
			return &l
		}
	}
	return nil
}

// bigger means less..
// Return 1 if original is bigger
// Return 0 if equal
// Return -1 if smaller
// Title always bigger than Item
func (l *Line) Compare(line Line) int {
	// Title
	if l.T != nil {
		if line.T == nil {
			return 1
		}
		if l.T.Level < line.T.Level {
			return 1
		} else if l.T.Level == line.T.Level {
			return 0
		}
		return -1
	}

	// Item
	if line.T != nil {
		return -1
	}
	if l.I.Level < line.I.Level {
		return 1
	} else if l.I.Level == line.I.Level {
		return 0
	}
	return -1

}

func (l *Line) GetKey() (string, bool) {
	if l.T != nil {
		return l.T.Name, true
	}

	return l.I.Key, false
}

func (l *Line) RFCRecord() bool {
	if l.I == nil {
		return false
	}

	if l.I.RFC != "" {
		return true
	}

	return false
}

func (l *Line) Debug(prefix string) {
	if l.T != nil {
		fmt.Printf("%s Title %s %d\n", prefix, l.T.Ref, l.T.Level)
	} else {
		fmt.Printf("%s Item %s %d\n", prefix, l.I.Value, l.I.Level)
	}
}

func GetLines(data []byte) []Line {
	stripRegs := []*regexp.Regexp{
		backtickReg,
	}

	str := string(data)
	for _, s := range stripRegs {
		str = s.ReplaceAllString(str, "")
	}
	// Hotfix, the mandatory field sometimes is useless
	str = strings.Replace(str, "REQUIRED)", "Required)", -1)

	var lines []Line
	strs := strings.Split(str, "\n")
	for _, s := range strs {
		l := NewLine(s)
		if l != nil {
			lines = append(lines, *l)
		}
	}
	return lines
}
