package feeds

// rss support
// validation done according to spec here:
//    http://cyber.law.harvard.edu/rss/rss.html

import (
	"encoding/xml"
	"fmt"
	"time"
)

// private wrapper around the RssFeed which gives us the <rss>..</rss> xml
type RssFeedXml struct {
	XMLName             xml.Name `xml:"rss"`
	Version             string   `xml:"version,attr"`
	ContentNamespace    string   `xml:"xmlns:content,attr"`
	DublinCoreNamespace string   `xml:"xmlns:dc,attr"`
	MediaNamespace      string   `xml:"xmlns:media,attr"`
	AtomNamespace       string   `xml:"xmlns:atom,attr"`
	Channel             *RssFeed
}

type RssContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Content string   `xml:",cdata"`
}

type RssImage struct {
	XMLName xml.Name `xml:"image"`
	Url     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Width   int      `xml:"width,omitempty"`
	Height  int      `xml:"height,omitempty"`
}

// deprecated
// will be removed soon
type RssTextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}

type RssFeed struct {
	XMLName        xml.Name  `xml:"channel"`
	Title          string    `xml:"title"`       // required
	Link           string    `xml:"link"`        // required
	Description    string    `xml:"description"` // required
	Language       string    `xml:"language,omitempty"`
	Copyright      string    `xml:"copyright,omitempty"`
	ManagingEditor string    `xml:"managingEditor,omitempty"` // Author used
	WebMaster      string    `xml:"webMaster,omitempty"`
	PubDate        string    `xml:"pubDate,omitempty"`       // created or updated
	LastBuildDate  string    `xml:"lastBuildDate,omitempty"` // updated used
	Category       string    `xml:"category,omitempty"`
	Generator      string    `xml:"generator,omitempty"`
	Docs           string    `xml:"docs,omitempty"`
	Cloud          string    `xml:"cloud,omitempty"`
	Ttl            int       `xml:"ttl,omitempty"`
	Rating         string    `xml:"rating,omitempty"`
	SkipHours      string    `xml:"skipHours,omitempty"`
	SkipDays       string    `xml:"skipDays,omitempty"`
	SelfLink       *NamespacedAtomLink `xml:"atom:link,omitempty"`
	Image          *RssImage
	TextInput      *RssTextInput
	Items          []*RssItem `xml:"item"`
}

type NamespacedAtomLink struct {
	XMLName xml.Name `xml:"atom:link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Length  string   `xml:"length,attr,omitempty"`
}

type RssItem struct {
	XMLName          xml.Name `xml:"item"`
	Title            string   `xml:"title"` // required
	MediaTitle       string   `xml:"media:title,omitempty"`
	Link             string   `xml:"link"`        // required
	Description      string   `xml:"description"` // required
	MediaDescription string   `xml:"media:description"`
	Content          *RssContent
	Author           string   `xml:"author,omitempty"`
	Category         []string `xml:"category,omitempty"`
	//	MediaCategory []string `xml:"media:category,omitempty"` // TODO implement correctly
	Comments       string        `xml:"comments,omitempty"`
	MediaContent   *MediaContent `xml:"media:content,omitempty"`
	Enclosure      *RssEnclosure
	Guid           string          `xml:"guid,omitempty"`    // Id used
	PubDate        string          `xml:"pubDate,omitempty"` // created or updated
	Source         string          `xml:"source,omitempty"`
	Creator        string          `xml:"dc:creator,omitempty"`
	MediaThumbnail *MediaThumbnail `xml:"media:thumbnail,omitempty"`
	MediaCopyright string          `xml:"media:copyright,omitempty"`
}

type RssEnclosure struct {
	//RSS 2.0 <enclosure url="http://example.com/file.mp3" length="123456789" type="audio/mpeg" />
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

type MediaContent struct {
	Url          string `xml:"url,attr"`
	FileSize     string `xml:"file_size,attr,omitempty"`
	Type         string `xml:"type,attr,omitempty"`
	Medium       string `xml:"medium,attr,omitempty"`
	IsDefault    string `xml:"isDefault,attr,omitempty"`
	Expression   string `xml:"expression,attr,omitempty"`
	Bitrate      string `xml:"bitrate,attr,omitempty"`
	Framerate    string `xml:"framerate,attr,omitempty"`
	Samplingrate string `xml:"samplingrate,attr,omitempty"`
	Channels     string `xml:"channels,attr,omitempty"`
	Duration     string `xml:"duration,attr,omitempty"`
	Height       string `xml:"height,attr,omitempty"`
	Width        string `xml:"width,attr,omitempty"`
	Lang         string `xml:"lang,attr,omitempty"`
}

type MediaThumbnail struct {
	Url    string `xml:"url,attr"`
	Height string `xml:"height,attr,omitempty"`
	With   string `xml:"with,attr,omitempty"`
	Time   string `xml:"time,attr,omitempty"`
}

type Rss struct {
	*Feed
}

// create a new RssItem with a generic Item struct's data
func newRssItem(i *Item) *RssItem {
	item := &RssItem{
		Title:       i.Title,
		Link:        i.Link.Href,
		Description: i.Description,
		Guid:        i.Id,
		PubDate:     anyTimeFormat(time.RFC1123Z, i.Created, i.Updated),
	}
	if len(i.Content) > 0 {
		item.Content = &RssContent{Content: i.Content}
	}
	if i.Source != nil {
		item.Source = i.Source.Href
	}

	// Define a closure
	if i.Enclosure != nil && i.Enclosure.Type != "" && i.Enclosure.Length != "" {
		item.Enclosure = &RssEnclosure{Url: i.Enclosure.Url, Type: i.Enclosure.Type, Length: i.Enclosure.Length}
	}

	if i.Author != nil {
		item.Author = i.Author.Name
	}
	return item
}

// create a new RssFeed with a generic Feed struct's data
func (r *Rss) RssFeed() *RssFeed {
	pub := anyTimeFormat(time.RFC1123Z, r.Created, r.Updated)
	build := anyTimeFormat(time.RFC1123Z, r.Updated)
	author := ""
	if r.Author != nil {
		author = r.Author.Email
		if len(r.Author.Name) > 0 {
			author = fmt.Sprintf("%s (%s)", r.Author.Email, r.Author.Name)
		}
	}

	var image *RssImage
	if r.Image != nil {
		image = &RssImage{Url: r.Image.Url, Title: r.Image.Title, Link: r.Image.Link, Width: r.Image.Width, Height: r.Image.Height}
	}

	channel := &RssFeed{
		Title:          r.Title,
		Link:           r.Link.Href,
		Description:    r.Description,
		ManagingEditor: author,
		PubDate:        pub,
		LastBuildDate:  build,
		Copyright:      r.Copyright,
		Image:          image,
	}
	for _, i := range r.Items {
		channel.Items = append(channel.Items, newRssItem(i))
	}
	return channel
}

// FeedXml returns an XML-Ready object for an Rss object
func (r *Rss) FeedXml() interface{} {
	// only generate version 2.0 feeds for now
	return r.RssFeed().FeedXml()

}

// FeedXml returns an XML-ready object for an RssFeed object
func (r *RssFeed) FeedXml() interface{} {
	return &RssFeedXml{
		Version:             "2.0",
		Channel:             r,
		ContentNamespace:    "http://purl.org/rss/1.0/modules/content/",
		DublinCoreNamespace: "http://purl.org/dc/elements/1.1/",
		MediaNamespace:      "http://search.yahoo.com/mrss/",
		AtomNamespace:       "http://www.w3.org/2005/Atom",
	}
}
