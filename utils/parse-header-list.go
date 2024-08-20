package utils

import (
	"strings"

	. "github.com/ncpa0cpl/ezs"
)

func ParseHeaderList(raw string) *Array[string] {
	elems := NewArray(strings.Split(raw, ";"))
	elems = elems.Filter(func(element string, _ int) bool {
		return element != ""
	})
	return elems
}
