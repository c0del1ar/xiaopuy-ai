package ingest

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"

	"github.com/c0del1ar/xiaopuy-ai/internal/rag"
	"golang.org/x/net/html"
)

type ParsedPage struct {
	Title       string
	Description string
	Canonical   string
	Content     string
}

func ParseHTML(raw, pageURL string) (ParsedPage, error) {
	if strings.TrimSpace(raw) == "" { return ParsedPage{}, fmt.Errorf("HTML is empty") }
	doc, err := html.Parse(strings.NewReader(raw)); if err != nil { return ParsedPage{}, fmt.Errorf("parse HTML: %w", err) }
	result := ParsedPage{}
	var text []string
	var walk func(*html.Node, bool)
	walk = func(n *html.Node, ignored bool) {
		if n.Type == html.ElementNode {
			tag := strings.ToLower(n.Data)
			if tag == "script" || tag == "style" || tag == "noscript" || tag == "template" || tag == "svg" || tag == "nav" || tag == "footer" { return }
			if tag == "title" { result.Title = nodeText(n); return }
			if tag == "meta" { name, content := attr(n,"name"), attr(n,"content"); if strings.EqualFold(name,"description") { result.Description = strings.TrimSpace(content) }; return }
			if tag == "link" && strings.EqualFold(attr(n,"rel"),"canonical") { result.Canonical = resolveURL(pageURL, attr(n,"href")); return }
			if tag == "header" { return }
			if tag == "body" { ignored = false }
		}
		if n.Type == html.TextNode && !ignored { if value := cleanText(n.Data); value != "" { text = append(text, value) } }
		for child := n.FirstChild; child != nil; child = child.NextSibling { walk(child, ignored) }
	}
	walk(doc, false)
	result.Content = strings.Join(text, "\n")
	if len([]rune(result.Content)) < 20 { return ParsedPage{}, fmt.Errorf("page has insufficient readable content") }
	return result, nil
}

func (p ParsedPage) Document(id, pageURL string) rag.Document {
	canonical := p.Canonical; if canonical == "" { canonical = pageURL }
	return rag.Document{ID:id, Source:"website", URL:canonical, Title:p.Title, Type:"webpage", Trust:"first_party", Content:p.Content, ContentHash:rag.ContentHash(p.Content)}
}

func nodeText(n *html.Node) string { var parts []string; var walk func(*html.Node); walk=func(x *html.Node){ if x.Type==html.TextNode { if v:=cleanText(x.Data); v!="" {parts=append(parts,v)} }; for c:=x.FirstChild;c!=nil;c=c.NextSibling{walk(c)} }; walk(n); return strings.Join(parts," ") }
func attr(n *html.Node, key string) string { for _, a := range n.Attr { if strings.EqualFold(a.Key,key) { return a.Val } }; return "" }
func resolveURL(base, ref string) string { if ref=="" {return ""}; b,err:=url.Parse(base); if err!=nil{return ref}; u,err:=url.Parse(ref); if err!=nil{return ref}; return b.ResolveReference(u).String() }
func cleanText(s string) string { return strings.Join(strings.FieldsFunc(s, func(r rune) bool{return unicode.IsSpace(r)}), " ") }
