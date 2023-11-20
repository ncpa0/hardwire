package views

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/ncpa0/htmx-framework/configuration"
	"github.com/ncpa0/htmx-framework/utils"
	"golang.org/x/net/html"
)

type PageView struct {
	root          string
	title         string
	filepath      string
	routePathname string
	document      *NodeProxy
	queryCache    map[string]*NodeProxy
	queryAllCache map[string][]*NodeProxy
}

type NodeProxy struct {
	parentRoot *PageView
	node       *html.Node
	raw        string
	etag       string
}

func nodeToString(node *html.Node) (string, error) {
	var b bytes.Buffer
	err := html.Render(&b, node)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func addClass(node *html.Node, class string) {
	var currentClassAttribute *html.Attribute
	for i, attr := range node.Attr {
		if attr.Key == "class" {
			currentClassAttribute = &node.Attr[i]
			break
		}
	}

	if currentClassAttribute == nil {
		node.Attr = append(node.Attr, html.Attribute{
			Key: "class",
			Val: class,
		})
	} else {
		currentClassAttribute.Val += " " + class
	}
}

func NewPageView(root string, filepath string) (*PageView, error) {
	doc, err := htmlquery.LoadDoc(path.Join(root, filepath))

	if err != nil {
		return nil, err
	}

	var rawHtml string
	var title string
	var routePathname string = filepath
	rawHtml, err = nodeToString(doc)
	if err != nil {
		return nil, err
	}

	titleNode := htmlquery.FindOne(doc, "//title")
	if titleNode != nil && titleNode.FirstChild != nil {
		title = titleNode.FirstChild.Data
	}

	if !path.IsAbs(routePathname) {
		routePathname = "/" + routePathname
	}
	if configuration.Current.StripExtension {
		routePathname = routePathname[:len(routePathname)-len(path.Ext(routePathname))]
	}

	hash := utils.Hash(rawHtml)

	return &PageView{
		root:          root,
		title:         title,
		filepath:      filepath,
		routePathname: routePathname,
		queryCache:    make(map[string]*NodeProxy),
		queryAllCache: make(map[string][]*NodeProxy),
		document: &NodeProxy{
			node: doc,
			raw:  rawHtml,
			etag: hash,
		},
	}, nil
}

func (v *PageView) GetFilepath() string {
	return v.filepath
}

func (v *PageView) GetRoutePathname() string {
	return v.routePathname
}

func pathSegments(p string) []string {
	segments := strings.Split(p, "/")
	n := 0
	for _, segment := range segments {
		if segment != "" {
			segments[n] = segment
			n++
		}
	}

	return segments[:n]
}

func (v *PageView) MatchesRoute(routePathname string) bool {
	fpathSegments := pathSegments(routePathname)
	viewPathSegments := pathSegments(v.routePathname)

	if len(fpathSegments) == len(viewPathSegments) {
		for i, seg := range fpathSegments {
			viewSeg := viewPathSegments[i]
			if seg != viewSeg && viewSeg[0] != ':' {
				return false
			}
		}
		return true
	}

	return false
}

func (v *PageView) QuerySelector(selector string) *utils.Option[NodeProxy] {
	// first check cache
	if cached, ok := v.queryCache[selector]; ok {
		return utils.NewOption(cached)
	}

	query := utils.NewTranslator(selector).XPathQuery()
	result := htmlquery.FindOne(v.document.node, query)

	if result == nil {
		return utils.Empty[NodeProxy]()
	}

	var b bytes.Buffer
	html.Render(&b, result)
	rawHtml := b.String()
	node := &NodeProxy{
		node:       result,
		raw:        rawHtml,
		etag:       utils.Hash(rawHtml),
		parentRoot: v,
	}

	v.queryCache[selector] = node

	return utils.NewOption(node)
}

func (v *PageView) QuerySelectorAll(selector string) []*NodeProxy {
	// first check cache
	if cached, ok := v.queryAllCache[selector]; ok {
		result := make([]*NodeProxy, len(cached))
		copy(result, cached)
		return result
	}

	query := utils.NewTranslator(selector).XPathQuery()
	nodeList := htmlquery.Find(v.document.node, query)

	result := make([]*NodeProxy, len(nodeList))
	for _, node := range nodeList {
		var b bytes.Buffer
		html.Render(&b, node)
		rawHtml := b.String()
		node := &NodeProxy{
			node:       node,
			raw:        rawHtml,
			etag:       utils.Hash(rawHtml),
			parentRoot: v,
		}
		result = append(result, node)
	}

	cacheEntry := make([]*NodeProxy, len(result))
	copy(cacheEntry, result)
	v.queryAllCache[selector] = cacheEntry

	return result
}

func (v *PageView) GetNode() *NodeProxy {
	return v.document
}

func (n *NodeProxy) ToHtml() string {
	if n.parentRoot != nil && n.parentRoot.title != "" {
		return fmt.Sprintf("<title>%s</title>\n%s", n.parentRoot.title, n.raw)
	}
	return n.raw
}

func (n *NodeProxy) GetEtag() string {
	return n.etag
}
