package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
)

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

func parseHnHtmlToRss(html string) *Rss {
	// Parses the Hacker News html and returns an rssfeed

	// Example HN link:
	/*
		<tr><td class="title" align="right" valign="top">30.</td><td><center><a id="up_5802043" href="https://news.ycombinator.com/vote?for=5802043&amp;dir=up&amp;whence=%62%65%73%74"><img src="hn_files/grayarrow.gif" border="0" hspace="2" vspace="3"></a><span id="down_5802043"></span></center></td><td class="title"><a href="http://www.bchanx.com/logos-in-pure-css-demo">Show HN: Logos in Pure CSS</a><span class="comhead"> (bchanx.com) </span></td></tr><tr><td colspan="2"></td><td class="subtext"><span id="score_5802043">267 points</span> by <a href="https://news.ycombinator.com/user?id=bchanx">bchanx</a> 4 days ago  | <a href="https://news.ycombinator.com/item?id=5802043">74 comments</a></td></tr>
	*/

	re := regexp.MustCompile(`<td class="title"><a href="(?P<url>[^"]*)">(?P<desc>[^<]*)</a>(?:<span class="comhead">(?P<domain>[^<]*))?.*?href="item\?id=(?P<id>\d+)"`)

	rssfeed := &Rss{Version: "2.0",
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
			rssfeed.Items = append(rssfeed.Items, RssItem{fmt.Sprintf("%s%s", desc, domain),
				url,
				comments,
				id,
				fmt.Sprintf("<p><a href=\"%s\">%s</a> %s<p/><p><a href=\"%s\">Comments</a></p>", url, desc, domain, comments)})
		}
	}
	return rssfeed
}

// hnRssHandler is a http handler that writes the HN RSS feed as the response
func hnRssHandler(w http.ResponseWriter, r *http.Request) {
	html, err := getHNsourceHttp()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	rssfeed := parseHnHtmlToRss(html)
	rssfeed.printXml(w)
}
