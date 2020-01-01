package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Rss is the internal representation for a RSS feed
type Rss struct {
	XMLName     xml.Name  `xml:"rss"`
	Version     string    `xml:"version,attr"`
	Title       string    `xml:"channel>title"`
	Link        string    `xml:"channel>link"`
	Description string    `xml:"channel>description"`
	Items       []RssItem `xml:"channel>item"`
}

// RssItem is the internal representation for a RSS item
type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Comments    string `xml:"comments"`
	Guid        string `xml:"guid"`
	Description string `xml:"description"`
}

// printXml serializes the RSS feed xml to an io.Writer
func (r *Rss) printXml(w io.Writer) error {
	// Writes the xml of the rssfeed to w
	fmt.Fprintln(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(r)
}
