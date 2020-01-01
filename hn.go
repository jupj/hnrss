package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// hnItem represents one link in the HN links list
type hnItem struct {
	url,
	desc,
	domain,
	guid,
	comments string
}

// helper functions for processing the html:

// getAttr retrieves attribute key from n. Returns empty string if not found
func getAttr(n *html.Node, key string) string {
	if n == nil {
		return ""
	}
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// hasAttr checks if n has an attribute key with value val
func hasAttr(n *html.Node, key, val string) bool {
	if n == nil {
		return false
	}
	for _, attr := range n.Attr {
		if attr.Key == key && attr.Val == val {
			return true
		}
	}
	return false
}

// isElement tests whether n is an element node of type elem
func isElement(n *html.Node, elem string) bool {
	return n != nil && n.Type == html.ElementNode && n.Data == elem
}

// parseHnHtmlToRss takes a reader for the HN html and returns an *Rss
// of the links found in the HN html
func parseHnHtmlToRss(r io.Reader) (rss *Rss, err error) {
	// Parses the Hacker News html and returns an rssfeed

	// Example HN link:
	/*
		<tr>
		    <td align=right valign=top class="title">1.</td>
		    <td><center><a id=up_6336178 href="vote?for=6336178&amp;dir=up&amp;whence=%62%65%73%74"><img src="grayarrow.gif" border=0 vspace=3 hspace=2 alt="upvote"></a><span id=down_6336178></span></center></td>
		    <td class="title">
		        <a href="http://www.nytimes.com/2013/09/06/us/nsa-foils-much-internet-encryption.html">N.S.A. Foils Much Internet Encryption</a>
		        <span class="comhead"> (nytimes.com) </span>
		    </td>
		</tr>
		<tr>
		    <td colspan=2></td>
		    <td class="subtext"><span id=score_6336178>882 points</span> by <a href="user?id=ebildsten">ebildsten</a> 5 days ago  | <a href="item?id=6336178">386 comments</a></td>
		</tr>
		<tr style="height:5px"></tr>
	*/
	doc, err := html.Parse(r)

	rss = &Rss{Version: "2.0",
		Title:       "Hacker News Top Links",
		Link:        "https://news.ycombinator.com/best",
		Description: "Links for the intellectually curious, ranked by readers."}

	// parse the html tree:
	item := hnItem{}
	var parse func(*html.Node)
	parse = func(n *html.Node) {
		// Find start of item
		if isElement(n, "td") && hasAttr(n, "class", "title") {
			if isElement(n.FirstChild, "a") {
				// Create a new link
				item.url = getAttr(n.FirstChild, "href")
				// Replace relative link to HN item with absolute link
				if strings.HasPrefix(item.url, "item?id=") {
					item.url = "https://news.ycombinator.com/" + item.url
				}
				item.desc = n.FirstChild.FirstChild.Data
				item.domain = ""
			}
		}

		// Find domain
		if isElement(n.Parent, "td") && hasAttr(n.Parent, "class", "title") &&
			isElement(n, "span") && hasAttr(n, "class", "comhead") {
			item.domain = n.FirstChild.Data
		}

		// Find comment - end of the item
		if isElement(n, "a") && strings.HasPrefix(getAttr(n, "href"), "item?id=") && strings.HasSuffix(n.FirstChild.Data, "comments") {
			item.guid = regexp.MustCompile(`item\?id=(?P<id>\d+)`).FindStringSubmatch(getAttr(n, "href"))[1]
			item.comments = fmt.Sprintf("https://news.ycombinator.com/item?id=%s", item.guid)

			// Reached the end of this item. Add to feed:
			rss.Items = append(rss.Items, RssItem{fmt.Sprintf("%s%s", item.desc, item.domain),
				item.url,
				item.comments,
				item.guid,
				fmt.Sprintf("<p><a href=\"%s\">%s</a> %s<p/><p><a href=\"%s\">Comments</a></p>", item.url, item.desc, item.domain, item.comments)})
			// reset item
			item = hnItem{}
		}

		// process the rest of the document
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			parse(c)
		}
	}
	parse(doc)

	return rss, err
}

// hnRssHandler is a http handler that writes the HN RSS feed as the response
func hnRssHandler(w http.ResponseWriter, r *http.Request) {
	// Read the Hacker News html from the web
	resp, err := http.Get("https://news.ycombinator.com/best")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	rssfeed, err := parseHnHtmlToRss(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	rssfeed.printXml(w)
}
