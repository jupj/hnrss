package main

import (
	"encoding/xml"
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
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Comments string `xml:"comments"`
    Guid string `xml:"guid"`
	Description string `xml:"description"`
}

// Print serializes the RSS feed to an io.Writer
func (r *Rss) print(w io.Writer) error {
	// Writes the xml of the rssfeed to w
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(r)
}
