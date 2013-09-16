package main

import (
	"fmt"
	//"html"
	"io"
	//"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"code.google.com/p/go.net/html"
)

type hnItem struct {
    url,
    desc,
    domain,
    guid,
    comments string
}

/*
func getHNsourceHttp() (string, error) {
	// Read the Hacker News html from the web
	resp, err := http.Get("https://news.ycombinator.com/best")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return html.UnescapeString(string(body)), nil
}
*/

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

    // helper functions for processing the html:
    getAttr := func(n *html.Node, key string) string {
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
    hasAttr := func(n *html.Node, key, val string) bool {
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
    isElement := func(n *html.Node, elem string) bool {
        return n != nil && n.Type == html.ElementNode && n.Data == elem
    }
    isHnLink := func(n *html.Node) bool {
        return isElement(n, "a") && isElement(n.Parent, "td") && hasAttr(n.Parent, "class", "title")
    }

    // parse the html tree:
    item := hnItem{}
    var f func(*html.Node)
    f = func(n *html.Node) {
        if isHnLink(n) {
            // Create a new link
            item.url = getAttr(n, "href")
            item.desc = n.FirstChild.Data
            item.domain = ""
            if isElement(n.NextSibling, "span") && getAttr(n.NextSibling, "class") == "comhead" {
                item.domain = n.NextSibling.FirstChild.Data
            }

            fmt.Println("url:", item.url)
            fmt.Println("desc:", item.desc)
            fmt.Println("domain:", item.domain)
        }

        if isElement(n, "a") && strings.HasPrefix(getAttr(n, "href"), "item?id=") && strings.HasSuffix(n.FirstChild.Data, "comments") {
            item.guid = regexp.MustCompile(`item\?id=(?P<id>\d+)`).FindStringSubmatch(getAttr(n, "href"))[1]
            item.comments = fmt.Sprintf("https://news.ycombinator.com/item?id=%s", item.guid)
            fmt.Println("guid:", item.guid)
            fmt.Println("comments:", item.comments)
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
            f(c)
        }
    }
    f(doc)

    /*
    htmlBytes, _ := ioutil.ReadAll(r)
    html := string(htmlBytes)

	re := regexp.MustCompile(`<td class="title"><a href="(?P<url>[^"]*)">(?P<desc>[^<]*)</a>(?:<span class="comhead">(?P<domain>[^<]*))?.*?href="item\?id=(?P<id>\d+)"`)

	rss := &Rss{Version: "2.0",
		Title:       "Hacker News Top Links",
		Link:        "https://news.ycombinator.com/best",
		Description: "Links for the intellectually curious, ranked by readers."}

	//matches := re.FindAllStringSubmatch(html, -1)
	matches := re.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) > 0 {
			url := match[1]
			desc := match[2]
			domain := match[3]
			id := match[4]
			comments := fmt.Sprintf("https://news.ycombinator.com/item?id=%s", id)
			rss.Items = append(rss.Items, RssItem{fmt.Sprintf("%s%s", desc, domain),
				url,
				comments,
				id,
				fmt.Sprintf("<p><a href=\"%s\">%s</a> %s<p/><p><a href=\"%s\">Comments</a></p>", url, desc, domain, comments)})
		}
	}
	*/
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
