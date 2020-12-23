package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Url  string
	Text string
}

func Links(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	linkMap := map[*html.Node]*Link{}
	traverse(doc, linkMap)
	links := []Link{}
	for _, v := range linkMap {
		links = append(links, *v)
	}
	return links, nil
}

func traverse(node *html.Node, linksMap map[*html.Node]*Link) {
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key != "href" {
					continue
				}
				var link *Link
				var ok bool
				if link, ok = linksMap[n.Parent]; !ok {
					link = &Link{attr.Val, ""}
					linksMap[n.Parent] = link
				}
				s := link.Text + " " + text(n)
				link.Text = strings.Join(strings.Fields(s), " ")
			}
		}
		traverse(n, linksMap)
	}
	return
}

func text(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	s := ""
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		s = s + " " + text(n)
	}
	return s
}

