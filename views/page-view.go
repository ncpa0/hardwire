package views

import (
	"bytes"
	"encoding/xml"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/antchfx/xmlquery"
	echo "github.com/labstack/echo/v4"
	"github.com/ncpa0/hardwire/configuration"
	. "github.com/ncpa0/hardwire/hw-context"

	// resources "github.com/ncpa0/hardwire/resource-provider"
	"github.com/ncpa0/hardwire/utils"
	. "github.com/ncpa0cpl/ezs"
)

type PageView struct {
	root              string
	title             string
	filepath          string
	routePathname     string
	isDynamic         bool
	requiredResources *Map[string, string]
	queryCache        *Map[string, *NodeProxy]
	queryAllCache     *Map[string, *Array[*NodeProxy]]
	document          *NodeProxy
	head              *NodeProxy
	Metadata          *pageMetafile
}

type NodeProxy struct {
	parentRoot *PageView
	node       *xmlquery.Node
	raw        string
	etag       string

	// Only present if parent's isDynamic is true
	template *template.Template
}

func addClass(node *xmlquery.Node, class string) {
	var currentClassAttribute *xmlquery.Attr
	for i, attr := range node.Attr {
		if attr.Name.Local == "class" {
			currentClassAttribute = &node.Attr[i]
			break
		}
	}

	if currentClassAttribute == nil {
		node.Attr = append(node.Attr, xmlquery.Attr{
			Name: xml.Name{
				Local: "class",
			},
			Value: class,
		})
	} else {
		currentClassAttribute.Value += " " + class
	}
}

func NewPageView(root string, filepath string) (*PageView, error) {
	file, err := os.Open(path.Join(root, filepath))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	doc, err := xmlquery.Parse(file)

	if err != nil {
		return nil, err
	}

	dirname := path.Dir(filepath)
	basename := path.Base(strings.TrimSuffix(filepath, ".html"))
	metaFilepath := path.Join(root, dirname, basename+".meta.json")

	metaFile, err := loadPageMetafile(metaFilepath)
	if err != nil {
		return nil, err
	}

	var rawHtml string
	var title string
	var routePathname string = filepath
	rawHtml = utils.XmlNodeToString(doc)

	titleNode := xmlquery.FindOne(doc, "//title")
	if titleNode != nil && titleNode.FirstChild != nil {
		title = titleNode.FirstChild.Data
	}

	if !path.IsAbs(routePathname) {
		routePathname = "/" + routePathname
	}
	if !configuration.Current.KeepExtension {
		routePathname = routePathname[:len(routePathname)-len(path.Ext(routePathname))]
	}

	hash := utils.Hash(rawHtml)

	var templ *template.Template
	if metaFile.IsDynamic {
		templ, err = template.New(filepath).Parse(rawHtml)

		if err != nil {
			return nil, err
		}
	}

	requiredResources := NewMap(map[string]string{})
	for _, res := range metaFile.Resources {
		requiredResources.Set(res.Key, res.Res)
	}

	view := &PageView{
		root:              root,
		title:             title,
		filepath:          filepath,
		routePathname:     routePathname,
		isDynamic:         metaFile.IsDynamic,
		requiredResources: requiredResources,
		queryCache:        NewMap(map[string]*NodeProxy{}),
		queryAllCache:     NewMap(map[string]*Array[*NodeProxy]{}),
		document: &NodeProxy{
			node:     doc,
			raw:      rawHtml,
			etag:     hash,
			template: templ,
		},
		Metadata: metaFile,
	}

	view.document.parentRoot = view

	headNode := xmlquery.FindOne(doc, "//head")
	if headNode != nil {
		view.head = &NodeProxy{
			parentRoot: view,
			node:       headNode,
			raw:        utils.XmlNodeToString(headNode),
		}
	}

	return view, nil
}

func (v *PageView) GetResourceKeys() *Array[string] {
	return v.requiredResources.Values()
}

func (v *PageView) IsDynamic() bool {
	return v.isDynamic
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
	if cached, ok := v.queryCache.Get(selector); ok {
		return utils.NewOption(cached)
	}

	query := utils.NewTranslator(selector).XPathQuery()
	result := xmlquery.FindOne(v.document.node, query)

	if result == nil {
		return utils.Empty[NodeProxy]()
	}

	rawHtml := utils.XmlNodeToString(result)

	var templ *template.Template
	if v.isDynamic {
		templ, _ = template.New(v.filepath + query).Parse(rawHtml)
	}

	node := &NodeProxy{
		node:       result,
		raw:        rawHtml,
		etag:       utils.Hash(rawHtml),
		parentRoot: v,
		template:   templ,
	}

	v.queryCache.Set(selector, node)

	return utils.NewOption(node)
}

func (v *PageView) QuerySelectorAll(selector string) []*NodeProxy {
	// first check cache
	if cached, ok := v.queryAllCache.Get(selector); ok {
		return cached.ToSlice()
	}

	query := utils.NewTranslator(selector).XPathQuery()
	nodeList := xmlquery.Find(v.document.node, query)

	result := make([]*NodeProxy, len(nodeList))
	for _, node := range nodeList {
		rawHtml := utils.XmlNodeToString(node)

		var templ *template.Template
		if v.isDynamic {
			templ, _ = template.New(v.filepath + query).Parse(rawHtml)
		}

		node := &NodeProxy{
			node:       node,
			raw:        rawHtml,
			etag:       utils.Hash(rawHtml),
			parentRoot: v,
			template:   templ,
		}
		result = append(result, node)
	}

	v.queryAllCache.Set(selector, NewArray(result).Copy())

	return result
}

type RenderedView struct {
	Html string
	Etag string
	Head string
}

func (v *PageView) Render(hw HardwireContext, c echo.Context) (*RenderedView, error) {
	return v.document.Render(hw, c)
}

func (node *NodeProxy) Render(hw HardwireContext, c echo.Context) (*RenderedView, error) {
	var rawHtml string
	var etag string

	if node.parentRoot.isDynamic {
		templateData := NewMap(map[string]interface{}{})

		for entry := range node.parentRoot.requiredResources.Iter() {
			handler, err := hw.GetResourceHandler(c, entry.Value)
			if err != nil {
				return nil, err
			}
			resourceValue, err := handler(node.parentRoot.routePathname, utils.ParamMap(c))
			if err != nil {
				return nil, err
			}

			templateData.Set(entry.Key, resourceValue)
		}

		var buff bytes.Buffer
		err := node.template.Execute(&buff, templateData.ToMap())
		if err != nil {
			return nil, err
		}

		rawHtml = buff.String()

		if configuration.Current.Caching.DynamicRoutes.NoStore {
			etag = ""
		} else {
			etag = utils.Hash(rawHtml)
		}
	} else {
		rawHtml = node.raw
		etag = node.etag
	}

	result := RenderedView{
		Html: rawHtml,
		Etag: etag,
	}

	if node.parentRoot.head != nil {
		result.Head = node.parentRoot.head.raw
	}

	return &result, nil
}
