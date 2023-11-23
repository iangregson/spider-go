package parser

import (
	"github.com/iangregson/spider-go/urls"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
	"io"
	"net/url"
)

func ParseLinks(htmlBody io.Reader, u *url.URL) []*url.URL {
	var urls []*url.URL

	doc, err := html.Parse(htmlBody)
	if err != nil {
		log.Error().Err(err).Msg("Couldn't parse html.")
		return urls
	}

	urls = extractLinks(doc, u)

	return urls
}

func extractLinks(n *html.Node, u *url.URL) []*url.URL {
	var links []*url.URL

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				isValid, link := urls.ValidUrl(attr.Val, u)
				if isValid {
					links = append(links, link)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, extractLinks(c, u)...)
	}

	return links
}
