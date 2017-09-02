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
type Leaf struct {
	// Type: Title, Item
	Type string
	// Value: original string
	Value string
	// Level
	Level int
	// Key
	// Title: it is the display name
	// Content: it is the **`$key`**, could be empty
	Key string
	// has MUST|SHOULD ..
	// Title: should be empty
	RFC string
}

func GetTitleNameAndRef(value string) (string, string) {
	name := strings.TrimSpace(value)

	ref := strings.Replace(name, "(", "", -1)
	ref = strings.Replace(ref, ")", "", -1)
	ref = strings.Replace(ref, " ", "-", -1)
	ref = strings.ToLower(ref)
	return name, ref
}

func NewTitle(value string) (Leaf, error) {
	var title Leaf
	// trim all the anchors
	src := anchorReg.ReplaceAllString(value, "")
	ts := titleReg.FindStringSubmatch(src)
	if len(ts) != 3 {
		return Leaf{}, fmt.Errorf("'%s' is not valid title", value)
	}

	title.Type = "title"
	title.Value = value
	title.Level = len(ts[1])
	title.Key = strings.TrimSpace(ts[2])
	return title, nil
}

func NewItem(value string) Leaf {
	var item Leaf

	item.Type = "item"
	item.Value = value
	// Item is the blank space numbers at the beginning
	//   the line which has more spaces, is the child most likely
	for _, v := range value {
		if v == ' ' {
			item.Level++
		} else {
			break
		}
	}

	// Item always lower than title
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

func NewLeaf(value string) *Leaf {
	if strings.HasPrefix(value, "#") {
		if title, err := NewTitle(value); err == nil {
			return &title
		}
	} else {
		item := NewItem(value)
		if item.Useful() {
			return &item
		}
	}
	return nil
}

func (i *Leaf) Useful() bool {
	if i.Type == "title" {
		return true
	}
	if i.Key != "" || i.RFC != "" {
		return true
	}

	return false
}

// bigger means 'parent'
// Return 1 if original is bigger
// Return 0 if equal
// Return -1 if smaller
// Title always bigger than Item
func (l *Leaf) Compare(line Leaf) int {
	if l.Level < line.Level {
		return 1
	} else if l.Level == line.Level {
		return 0
	}

	return -1
}

func (l *Leaf) GetLevel() int {
	return l.Level
}

func (l *Leaf) GetKey() string {
	return l.Key
}

func (l *Leaf) RFCRecord() bool {
	if l.RFC != "" {
		return true
	}

	return false
}

func (l *Leaf) Debug(prefix string) {
	fmt.Printf("%s %s %s %d: %s\n", prefix, l.Type, l.Key, l.Level, l.Value)
}

func GetLeafs(data []byte) []Leaf {
	stripRegs := []*regexp.Regexp{
		backtickReg,
	}

	str := string(data)
	for _, s := range stripRegs {
		str = s.ReplaceAllString(str, "")
	}
	// Hotfix, the mandatory field sometimes is useless
	str = strings.Replace(str, "REQUIRED)", "Required)", -1)

	var leafs []Leaf
	strs := strings.Split(str, "\n")
	for _, s := range strs {
		l := NewLeaf(s)
		if l != nil {
			leafs = append(leafs, *l)
		}
	}
	return leafs
}
