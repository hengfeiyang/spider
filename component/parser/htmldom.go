package parser

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// HTMLDom HTML Selector 解析器，类似jQuery的选择器, document.querySelector
type HTMLDom struct {
	dom *goquery.Document // HTML Dom结构
}

// NewHTMLDom 创建一个HTML Selector 解析器
func NewHTMLDom(body []byte) (*HTMLDom, error) {
	if body == nil {
		return nil, Errorf("body is empty")
	}
	var err error
	t := new(HTMLDom)
	t.dom, err = goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Match 配置规则，返回匹配到的值
func (t *HTMLDom) Match(rule string) string {
	dom := t.DomFind(rule)
	if dom != nil {
		v, _ := dom.Html()
		return v
	}
	return ""
}

// MatchAll 配置规则，返回匹配到的值，复数
func (t *HTMLDom) MatchAll(rule string) []string {
	var vals []string
	dom := t.DomFind(rule)
	for i := 0; i < dom.Size(); i++ {
		node := dom.Get(i)
		vals = append(vals, t.NodeHTML(node))
	}
	return vals
}

// DomFind use document find a selector return html
// you can use chrome to select the dom selector, then .Html() fetch the special dom html
func (t *HTMLDom) DomFind(selector string) *goquery.Selection {
	if selector == "" {
		return nil
	}
	return t.dom.Find(selector)
}

// NodeText return node text
func (t *HTMLDom) NodeText(node *html.Node) string {
	if node == nil {
		return ""
	}
	return t.GetNodeText(node)
}

// NodeHTML return node html
func (t *HTMLDom) NodeHTML(node *html.Node) string {
	if node == nil {
		return ""
	}
	// Since there is no .innerHtml, the HTML content must be re-created from
	// the nodes using html.Render.
	var buf bytes.Buffer
	var e error
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		e = html.Render(&buf, c)
		if e != nil {
			return ""
		}
	}
	return buf.String()
}

// NodeAttr return node attr value
func (t *HTMLDom) NodeAttr(node *html.Node, attr string) string {
	if node == nil || attr == "" {
		return ""
	}

	for i, a := range node.Attr {
		if a.Key == attr {
			return node.Attr[i].Val
		}
	}
	return ""
}

// GetNodeText Get the specified node's text content.
func (t *HTMLDom) GetNodeText(node *html.Node) string {
	if node.Type == html.TextNode {
		// Keep newlines and spaces, like jQuery
		return node.Data
	} else if node.FirstChild != nil {
		var buf bytes.Buffer
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			buf.WriteString(t.GetNodeText(c))
		}
		return buf.String()
	}

	return ""
}
