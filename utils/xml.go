package utils

import (
	"strings"

	"github.com/antchfx/xmlquery"
)

const PREFIX_LEN = len("<?xml version=\"1.0\"?>")
const TEMPL_QUOTE = "@#34T;"

func XmlNodeToString(node *xmlquery.Node) string {
	s := node.OutputXMLWithOptions(
		xmlquery.WithOutputSelf(),
		xmlquery.WithPreserveSpace(),
	)
	if strings.HasPrefix(s, "<?xml") {
		s = s[PREFIX_LEN:]
	}
	s = strings.ReplaceAll(s, TEMPL_QUOTE, "\"")
	return s
}
