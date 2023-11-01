package utils

import (
	"regexp"
	"strings"
)

const (
	EQUALS_EXACT                     = "="
	EQUALS_CONTAINS_WORD             = "~="
	EQUALS_ENDS_WITH                 = "$="
	EQUALS_CONTAINS                  = "*="
	EQUALS_STARTS_WITH_OR_HYPHENATED = "|="
	EQUALS_STARTS_WITH               = "^="
)

var CSS_REG = regexp.MustCompile(`(?P<star>\*)|(:(?P<pseudo>[\w-]*))|\(*(?P<pseudospecifier>["\']*[\w\s-]*["\']*)\)|(?P<child>\s*>\s*)|(#(?P<id>[\w-]*))|(\.(?P<class>[\w-]*))|(?P<sibling>\s*\+\s*)|(\[(?P<attribute>[\w-]*)((?P<attribute_equals>[=~$*]+)(?P<attribute_value>(.+\[\]'?)|[^\]]+))*\])+|(?P<descendant>\s+)|(?P<element>[\w-]*)`)

type Translator struct {
	cssSelector string
	prefix      string
}

func NewTranslator(cssSelector string) *Translator {
	return &Translator{
		cssSelector: cssSelector,
		prefix:      "//",
	}
}

func (t *Translator) SetPrefix(prefix string) *Translator {
	t.prefix = prefix
	return t
}

func (t *Translator) XPathQuery() string {
	return t.Convert(t.cssSelector)
}

const COMMA = byte(',')
const QUOTE = byte('\'')
const DOUBLE_QUOTE = byte('"')

func SplitCss(css string) []string {
	results := []string{""}
	isInQuote := false
	currentQuoteSym := byte('.')

	for i := 0; i < len(css); i++ {
		char := css[i]
		if char == QUOTE {
			if isInQuote && currentQuoteSym == QUOTE {
				isInQuote = false
			} else {
				isInQuote = true
				currentQuoteSym = QUOTE
			}
		} else if char == DOUBLE_QUOTE {
			if isInQuote && currentQuoteSym == DOUBLE_QUOTE {
				isInQuote = false
			} else {
				isInQuote = true
				currentQuoteSym = DOUBLE_QUOTE
			}
		} else if char == COMMA {
			if !isInQuote {
				results = append(results, "")
				continue
			}
		}

		results[len(results)-1] += string(char)
	}

	if results[len(results)-1] == "" {
		results = results[:len(results)-1]
	}

	return results
}

func (t *Translator) Convert(css string) string {
	cssArray := SplitCss(css)

	xPathArray := make([]string, 0)

	for _, input := range cssArray {
		output := t.convertSingleSelector(strings.TrimSpace(input))
		xPathArray = append(xPathArray, output)
	}

	return strings.Join(xPathArray, " | ")
}

func mapToArr[T any](m map[string]T) []T {
	arr := make([]T, len(m))

	i := 0
	for _, v := range m {
		arr[i] = v
		i++
	}

	return arr
}

func (y *Translator) convertSingleSelector(css string) string {
	thread := pregMatchCollated(CSS_REG, css)

	xpath := []string{y.prefix}
	hasElement := false

	for threadKey, currentThreadItem := range thread {
		var next *MatchDetail
		if threadKey+1 < len(thread) {
			next = thread[threadKey+1]
		}

		switch currentThreadItem.Type {
		case "star", "element":
			xpath = append(xpath, currentThreadItem.Content)
			hasElement = true
		case "pseudo":
			specifier := ""
			if next.Type == "pseudospecifier" {
				specifier = next.Content
			}

			switch currentThreadItem.Content {
			case "disabled", "checked", "selected":
				xpath = append(xpath, "[@"+currentThreadItem.Content+"]")
			case "text":
				xpath = append(xpath, `[@type="text"]`)
			case "contains":
				if specifier == "" {
					continue
				}
				xpath = append(xpath, "[contains(text(),"+specifier+")]")
			case "first-child":
				prev := len(xpath) - 1
				xpath[prev] = `*[1]/self::` + xpath[prev]
			case "nth-child":
				if specifier == "" {
					continue
				}
				prev := len(xpath) - 1
				previous := xpath[prev]
				if strings.HasSuffix(previous, "]") {
					xpath[prev] = strings.Replace(previous, "]", " and position() = "+specifier+"]", 1)
				} else {
					xpath = append(xpath, "["+specifier+"]")
				}
			case "nth-of-type":
				if specifier == "" {
					continue
				}
				prev := len(xpath) - 1
				previous := xpath[prev]
				if strings.HasSuffix(previous, "]") {
					xpath = append(xpath, "["+specifier+"]")
				} else {
					xpath = append(xpath, "["+specifier+"]")
				}
			}
		case "child":
			xpath = append(xpath, "/")
			hasElement = false
		case "id":
			xpath = append(xpath, (func() string {
				if hasElement {
					return ""
				}
				return "*"
			}())+"[@id='"+currentThreadItem.Content+"']")
			hasElement = true
		case "class":
			// https://devhints.io/xpath#class-check
			xpath = append(xpath, (func() string {
				if hasElement {
					return ""
				}
				return "*"
			}())+"[contains(concat(' ', normalize-space(@class), ' '), ' "+currentThreadItem.Content+" ')]")
			hasElement = true
		case "sibling":
			xpath = append(xpath, "/following-sibling::*[1]/self::")
			hasElement = false
		case "attribute":
			if !hasElement {
				xpath = append(xpath, "*")
				hasElement = true
			}

			if len(currentThreadItem.Detail) == 0 {
				xpath = append(xpath, "[@"+currentThreadItem.Content+"]")
				continue
			}

			cmpType := currentThreadItem.Detail[0]
			cmpValue := currentThreadItem.Detail[1]

			valueString := strings.Trim(cmpValue.Content, " '\"")
			equalsType := cmpType.Content

			switch equalsType {
			case EQUALS_EXACT:
				xpath = append(xpath, "[@"+currentThreadItem.Content+"='"+valueString+"']")
			case EQUALS_CONTAINS:
				xpath = append(xpath, "[contains(@"+currentThreadItem.Content+",'"+valueString+"')]")
			case EQUALS_CONTAINS_WORD:
				xpath = append(xpath, "[contains(concat(' ',@"+currentThreadItem.Content+",' '),' "+valueString+" ')]")
			case EQUALS_STARTS_WITH:
				panic("Not Yet Implemented")
			case EQUALS_STARTS_WITH_OR_HYPHENATED:
				panic("Not Yet Implemented")
			case EQUALS_ENDS_WITH:
				xpath = append(xpath, "[substring(@"+currentThreadItem.Content+",string-length(@"+currentThreadItem.Content+") - string-length('"+valueString+"') + 1)='"+valueString+"']")
			}
		case "descendant":
			xpath = append(xpath, "//")
			hasElement = false
		}
	}

	return strings.Join(xpath, "")
}

type MatchDetail struct {
	Type    string
	Content string
	Detail  []MatchDetail
}

type MatchSet map[string]MatchDetail

type Tuple[T any, U any] struct {
	Value1 T
	Value2 U
}

func extractNamedCaptureGroups(match []string, subexpNames []string) []Tuple[string, string] {
	result := []Tuple[string, string]{}
	for i, name := range subexpNames {
		if i != 0 && name != "" {
			value := match[i]

			if value != "" {
				result = append(result, Tuple[string, string]{name, match[i]})
			}
		}
	}
	return result
}

func pregMatchCollated(re *regexp.Regexp, str string) []*MatchDetail {
	// set := map[string]MatchDetail{}
	result := []*MatchDetail{}

	// add := func(t string, content string) *MatchDetail {
	// 	det := MatchDetail{
	// 		Type:    t,
	// 		Content: content,
	// 	}

	// 	// if _, ok := set[t]; ok {
	// 	// 	mainDet := set[t]
	// 	// 	mainDet.Detail = append(mainDet.Detail, det)
	// 	// } else {
	// 	// 	set[t] = det
	// 	// }
	// 	result = append(result, det)

	// 	return &det
	// }

	if matches := re.FindAllStringSubmatch(str, -1); matches != nil {
		for _, match := range matches {
			captureGroups := extractNamedCaptureGroups(match, re.SubexpNames())
			var first *MatchDetail
			for _, tuple := range captureGroups {
				name := tuple.Value1
				value := tuple.Value2

				if first != nil {
					first.Detail = append(first.Detail, MatchDetail{
						Type:    name,
						Content: value,
					})
				} else {
					first = &MatchDetail{
						Type:    name,
						Content: value,
						Detail:  []MatchDetail{},
					}
					result = append(result, first)
				}

				// fmt.Printf("%s: %s\n", name, value)
				// first = add(name, value)
			}
		}
	}

	return result
}
