package views

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/ncpa0/htmx-framework/utils"
	"golang.org/x/net/html"
)

type View struct {
	root     string
	filepath string
	raw      string
	document *html.Node
}

func NewView(root string, filepath string) (*View, error) {
	doc, err := htmlquery.LoadDoc(path.Join(root, filepath))

	if err != nil {
		fmt.Println("Error loading view.")
		fmt.Println(err)
		return nil, err
	}

	var b bytes.Buffer
	err = html.Render(&b, doc)

	if err != nil {
		return nil, err
	}

	return &View{
		root:     root,
		filepath: filepath,
		raw:      b.String(),
		document: doc,
	}, nil
}

func (v *View) GetFilepath() string {
	return v.filepath
}

func (v *View) FilepathMatches(fpath string) bool {
	if v.filepath == fpath {
		return true
	}

	fullPath := path.Join(v.root, v.filepath)
	return fullPath == fpath || strings.HasSuffix(fullPath, fpath)
}

func (v *View) QuerySelector(selector string) *utils.Option[html.Node] {
	query := utils.NewTranslator(selector).XPathQuery()
	result := htmlquery.FindOne(v.document, query)
	return utils.NewOption(result)
}

func (v *View) QuerySelectorAll(selector string) *utils.Option[[]*html.Node] {
	query := utils.NewTranslator(selector).XPathQuery()
	result := htmlquery.Find(v.document, query)
	return utils.NewOption(&result)
}

func (v *View) ToNode() *html.Node {
	return v.document
}
