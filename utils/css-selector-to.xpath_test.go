package utils_test

import (
	"testing"

	"github.com/ncpa0/hardwire/utils"
	"github.com/stretchr/testify/assert"
)

func TestCssSplit(t *testing.T) {
	ass := assert.New(t)

	ass.Equal([]string{"foo"}, utils.SplitCss("foo"))
	ass.Equal([]string{"a", "b"}, utils.SplitCss("a,b"))
	ass.Equal([]string{"a", "b\"c,d\"e"}, utils.SplitCss("a,b\"c,d\"e"))
	ass.Equal([]string{"',a,b,'", "'def'"}, utils.SplitCss("',a,b,','def'"))
}

func TestTranslator(t *testing.T) {
	ass := assert.New(t)

	assertSelector := func(selector string, query string) {
		result := utils.NewTranslator(selector).XPathQuery()
		ass.Equal(query, result)
		t.Logf("PASS: `%s` -> `%s`", selector, query)
	}

	assertSelector("#root", `//*[@id='root']`)
	assertSelector("button", "//button")
	assertSelector(".foo.bar.baz", "//*[contains(concat(' ', normalize-space(@class), ' '), ' foo ')][contains(concat(' ', normalize-space(@class), ' '), ' bar ')][contains(concat(' ', normalize-space(@class), ' '), ' baz ')]")
	assertSelector("div.container", "//div[contains(concat(' ', normalize-space(@class), ' '), ' container ')]")
	assertSelector("div#root > .elem[value='0']", "//div[@id='root']/*[contains(concat(' ', normalize-space(@class), ' '), ' elem ')][@value='0']")
}
