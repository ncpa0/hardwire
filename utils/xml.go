package utils

import (
	"encoding/xml"
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

func XmlNodeClone(node *xmlquery.Node) *xmlquery.Node {
	deref := *node
	nodeCopy := deref
	nodeCopy.Attr = make([]xmlquery.Attr, len(node.Attr))
	copy(nodeCopy.Attr, node.Attr)
	return &nodeCopy
}

func XmlNodeSetAttribute(node *xmlquery.Node, attribute string, value string) {
	node.Attr = append(node.Attr, xmlquery.Attr{
		Name: xml.Name{
			Local: attribute,
		},
		Value: value,
	})
}
