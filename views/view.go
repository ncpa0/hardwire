package views

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/antchfx/htmlquery"
	"github.com/ncpa0/htmx-framework/configuration"
	"github.com/ncpa0/htmx-framework/utils"
	"golang.org/x/net/html"
)

type View struct {
	IsDynamicFragment bool
	Template          *template.Template
	RequiredResource  string
	root              string
	title             string
	filepath          string
	routePathname     string
	document          *NodeProxy
	queryCache        map[string]*NodeProxy
	queryAllCache     map[string][]*NodeProxy
}

type NodeProxy struct {
	parentRoot *View
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

func NewView(root string, filepath string) (*View, error) {
	doc, err := htmlquery.LoadDoc(path.Join(root, filepath))

	if err != nil {
		return nil, err
	}

	isDynamicFragment := false
	var viewTemplate *template.Template
	var templateRequiredResource string
	var rawHtml string
	var title string
	var routePathname string = filepath
	if strings.HasSuffix(filepath, ".template.html") {
		dirname := path.Dir(filepath)
		basename := path.Base(strings.TrimSuffix(filepath, ".template.html"))
		metaFilepath := path.Join(root, dirname, basename+".meta.json")

		metaFile, err := loadMetafile(metaFilepath)
		if err != nil {
			return nil, err
		}

		dynamicFragment := htmlquery.FindOne(doc, "//dynamic-fragment")

		if dynamicFragment == nil {
			return nil, errors.New("given template is not a valid dynamic fragment")
		}

		dynamicFragment.Data = "div"
		dynamicFragment.Attr = append(dynamicFragment.Attr, html.Attribute{
			Key: "data-frag-url",
			Val: filepath[:len(filepath)-len(".template.html")],
		})
		addClass(dynamicFragment, "__dynamic_fragment")
		rawHtml, err = nodeToString(dynamicFragment)

		if err != nil {
			return nil, err
		}

		templ, err := template.New(filepath).Parse(rawHtml)

		if err != nil {
			return nil, err
		}

		if !path.IsAbs(routePathname) {
			routePathname = "/" + routePathname
		}
		routePathname = routePathname[:len(routePathname)-len(".template.html")]

		isDynamicFragment = true
		viewTemplate = templ
		templateRequiredResource = metaFile.ResourceName
		doc = nil
	} else {
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
	}

	hash := utils.Hash(rawHtml)

	return &View{
		root:              root,
		title:             title,
		filepath:          filepath,
		routePathname:     routePathname,
		queryCache:        make(map[string]*NodeProxy),
		queryAllCache:     make(map[string][]*NodeProxy),
		IsDynamicFragment: isDynamicFragment,
		Template:          viewTemplate,
		RequiredResource:  templateRequiredResource,
		document: &NodeProxy{
			node: doc,
			raw:  rawHtml,
			etag: hash,
		},
	}, nil
}

func (v *View) GetFilepath() string {
	return v.filepath
}

func (v *View) GetRoutePathname() string {
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

func (v *View) MatchesRoute(routePathname string) bool {
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

func (v *View) QuerySelector(selector string) *utils.Option[NodeProxy] {
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

func (v *View) QuerySelectorAll(selector string) []*NodeProxy {
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

func (v *View) GetNode() *NodeProxy {
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

type TemplateMetafile struct {
	ResourceName string `json:"resourceName"`
	Hash         string `json:"hash"`
}

func loadMetafile(filepath string) (*TemplateMetafile, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metafile TemplateMetafile
	err = json.NewDecoder(file).Decode(&metafile)
	if err != nil {
		return nil, err
	}

	return &metafile, nil
}
