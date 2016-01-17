package templates

import (
	"html/template"
	"unicode/utf8"
)

var funcMap = template.FuncMap{
	"sum":          sum,
	"NewFirstChar": NewFirstChar,
}

func sum(a ...int) int {
	s := 0
	for _, x := range a {
		s += x
	}
	return s
}

// FirstChar type
type FirstChar struct {
	char rune
}

// NewFirstChar initialize new FirstChar
func NewFirstChar() *FirstChar {
	return &FirstChar{}
}

// NotEqual function
func (fc *FirstChar) NotEqual(str string) bool {

	if len(str) == 0 {
		fc.char = 0
		return true
	}

	r, _ := utf8.DecodeRuneInString(str)

	if fc.char == r {
		return false
	}
	fc.char = r
	return true
}

// String function
func (fc *FirstChar) String() string {
	return string(fc.char)
}
