package views

import (
	"bytes"
	"encoding/xml"
	"errors"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/antchfx/xmlquery"
)

type DynamicFragmentView struct {
	id               string
	template         *template.Template
	requiredResource string
	filepath         string
	routePathname    string
}

func NewDynamicFragmentView(root string, filepath string) (*DynamicFragmentView, error) {
	docFile, err := os.Open(path.Join(root, filepath))
	if err != nil {
		return nil, err
	}
	defer docFile.Close()

	doc, err := xmlquery.Parse(docFile)

	if err != nil {
		return nil, err
	}

	routePathname := filepath
	dirname := path.Dir(filepath)
	basename := path.Base(strings.TrimSuffix(filepath, ".template.html"))
	metaFilepath := path.Join(root, dirname, basename+".meta.json")

	metaFile, err := loadFragmentMetafile(metaFilepath)
	if err != nil {
		return nil, err
	}

	dynamicFragment := xmlquery.FindOne(doc, "//dynamic-fragment")

	if dynamicFragment == nil {
		return nil, errors.New("given template is not a valid dynamic fragment")
	}

	dynamicFragment.Data = "div"
	dynamicFragment.Attr = append(dynamicFragment.Attr, xmlquery.Attr{
		Name: xml.Name{
			Local: "data-frag-url",
		},
		Value: filepath[:len(filepath)-len(".template.html")],
	})
	addClass(dynamicFragment, "__dynamic_fragment")
	rawHtml := nodeToString(dynamicFragment)

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

	return &DynamicFragmentView{
		id:               metaFile.Hash,
		template:         templ,
		requiredResource: metaFile.ResourceName,
		filepath:         filepath,
		routePathname:    routePathname,
	}, nil
}

func (v *DynamicFragmentView) GetRoutePathname() string {
	return v.routePathname
}

func (v *DynamicFragmentView) GetResourceName() string {
	return v.requiredResource
}

func (v *DynamicFragmentView) MatchesRoute(routePathname string) bool {
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

func (v *DynamicFragmentView) Build(resource interface{}) (string, error) {
	var buff bytes.Buffer
	err := v.template.Execute(&buff, resource)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}
