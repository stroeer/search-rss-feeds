package feeds

import (
	"bytes"
	"encoding/xml"
	"testing"
	"time"
)

var rssOutput = `<channel>
  <title>jmoiron.net blog</title>
  <link>http://jmoiron.net/blog</link>
  <description>discussion about tech, footie, photos</description>
  <copyright>This work is copyright © Benjamin Button</copyright>
  <item>
    <title>Limiting Concurrency in Go</title>
    <link>http://jmoiron.net/blog/limiting-concurrency-in-go/</link>
    <description>A discussion on controlled parallelism in golang</description>
    <content:encoded><![CDATA[<p>Go's goroutines make it easy to make <a href="http://collectiveidea.com/blog/archives/2012/12/03/playing-with-go-embarrassingly-parallel-scripts/">embarrassingly parallel programs</a>, but in many &quot;real world&quot; cases resources can be limited and attempting to do everything at once can exhaust your access to them.</p>]]></content:encoded>
    <author>Jason Moiron</author>
    <guid isPermalink="true">http://jmoiron.net/blog/limiting-concurrency-in-go/</guid>
  </item>
  <item>
    <title>Limiting Concurrency in Go</title>
    <link>http://jmoiron.net/blog/limiting-concurrency-in-go/</link>
    <description>A discussion on controlled parallelism in golang</description>
    <content:encoded><![CDATA[<p>Go's goroutines make it easy to make <a href="http://collectiveidea.com/blog/archives/2012/12/03/playing-with-go-embarrassingly-parallel-scripts/">embarrassingly parallel programs</a>, but in many &quot;real world&quot; cases resources can be limited and attempting to do everything at once can exhaust your access to them.</p>]]></content:encoded>
    <author>Jason Moiron</author>
    <guid isPermalink="false">123456789</guid>
  </item>
  <item>
    <title>Logic-less Template Redux</title>
    <link>http://jmoiron.net/blog/logicless-template-redux/</link>
    <description>More thoughts on logicless templates</description>
  </item>
  <item>
    <title>Idiomatic Code Reuse in Go</title>
    <link>http://jmoiron.net/blog/idiomatic-code-reuse-in-go/</link>
    <description>How to use interfaces &lt;em&gt;effectively&lt;/em&gt;</description>
    <enclosure url="http://example.com/cover.jpg" length="123456" type="image/jpg"></enclosure>
  </item>
  <item>
    <title>Never Gonna Give You Up Mp3</title>
    <link>http://example.com/RickRoll.mp3</link>
    <description>Never gonna give you up - Never gonna let you down.</description>
    <enclosure url="http://example.com/RickRoll.mp3" length="123456" type="audio/mpeg"></enclosure>
  </item>
  <item>
    <title>String formatting in Go</title>
    <link>http://example.com/strings</link>
    <description>How to use things like %s, %v, %d, etc.</description>
  </item>
</channel>`

func TestRSS(t *testing.T) {
	feed := &RssFeed{
		Title:       "jmoiron.net blog",
		Link:        "http://jmoiron.net/blog",
		Description: "discussion about tech, footie, photos",
		Copyright:   "This work is copyright © Benjamin Button",
	}

	feed.Items = []*RssItem{
		{
			Title:       "Limiting Concurrency in Go",
			Link:        "http://jmoiron.net/blog/limiting-concurrency-in-go/",
			Description: "A discussion on controlled parallelism in golang",
			Author:      "Jason Moiron",
			Guid:        &RssGuid{IsPermalink: true, Guid: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
			Content:     &RssContent{Content: `<p>Go's goroutines make it easy to make <a href="http://collectiveidea.com/blog/archives/2012/12/03/playing-with-go-embarrassingly-parallel-scripts/">embarrassingly parallel programs</a>, but in many &quot;real world&quot; cases resources can be limited and attempting to do everything at once can exhaust your access to them.</p>`},
		},
		{
			Title:       "Limiting Concurrency in Go",
			Link:        "http://jmoiron.net/blog/limiting-concurrency-in-go/",
			Description: "A discussion on controlled parallelism in golang",
			Author:      "Jason Moiron",
			Guid:        &RssGuid{IsPermalink: false, Guid: "123456789"},
			Content:     &RssContent{Content: `<p>Go's goroutines make it easy to make <a href="http://collectiveidea.com/blog/archives/2012/12/03/playing-with-go-embarrassingly-parallel-scripts/">embarrassingly parallel programs</a>, but in many &quot;real world&quot; cases resources can be limited and attempting to do everything at once can exhaust your access to them.</p>`},
		},
		{
			Title:       "Logic-less Template Redux",
			Link:        "http://jmoiron.net/blog/logicless-template-redux/",
			Description: "More thoughts on logicless templates",
		},
		{
			Title:       "Idiomatic Code Reuse in Go",
			Link:        "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/",
			Description: "How to use interfaces <em>effectively</em>",
			Enclosure:   &RssEnclosure{Url: "http://example.com/cover.jpg", Length: "123456", Type: "image/jpg"},
		},
		{
			Title:       "Never Gonna Give You Up Mp3",
			Link:        "http://example.com/RickRoll.mp3",
			Enclosure:   &RssEnclosure{Url: "http://example.com/RickRoll.mp3", Length: "123456", Type: "audio/mpeg"},
			Description: "Never gonna give you up - Never gonna let you down.",
		},
		{
			Title:       "String formatting in Go",
			Link:        "http://example.com/strings",
			Description: "How to use things like %s, %v, %d, etc.",
		}}

	rss, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		t.Errorf("unexpected error encoding RSS: %v", err)
	}
	if string(rss) != rssOutput {
		t.Errorf("Rss not what was expected.  Got:\n%s\n\nExpected:\n%s\n", rss, rssOutput)
	}

	// if err := feed.WriteRss(&buf); err != nil {
	// 	t.Errorf("unexpected error writing RSS: %v", err)
	// }
	// if got := buf.String(); got != rssOutput {
	// 	t.Errorf("Rss not what was expected.  Got:\n%s\n\nExpected:\n%s\n", got, rssOutput)
	// }
}

var rssOutputSorted = `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:media="http://search.yahoo.com/mrss/" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>jmoiron.net blog</title>
    <link>http://jmoiron.net/blog</link>
    <description>discussion about tech, footie, photos</description>
    <copyright>This work is copyright © Benjamin Button</copyright>
    <managingEditor>jmoiron@jmoiron.net (Jason Moiron)</managingEditor>
    <pubDate>Wed, 16 Jan 2013 21:52:35 -0500</pubDate>
    <item>
      <title>Limiting Concurrency in Go</title>
      <link>http://jmoiron.net/blog/limiting-concurrency-in-go/</link>
      <description></description>
      <pubDate>Fri, 18 Jan 2013 21:52:35 -0500</pubDate>
    </item>
    <item>
      <title>Logic-less Template Redux</title>
      <link>http://jmoiron.net/blog/logicless-template-redux/</link>
      <description></description>
      <pubDate>Thu, 17 Jan 2013 21:52:35 -0500</pubDate>
    </item>
    <item>
      <title>Idiomatic Code Reuse in Go</title>
      <link>http://jmoiron.net/blog/idiomatic-code-reuse-in-go/</link>
      <description></description>
      <pubDate>Thu, 17 Jan 2013 09:52:35 -0500</pubDate>
    </item>
    <item>
      <title>Never Gonna Give You Up Mp3</title>
      <link>http://example.com/RickRoll.mp3</link>
      <description></description>
      <pubDate>Thu, 17 Jan 2013 07:52:35 -0500</pubDate>
    </item>
    <item>
      <title>String formatting in Go</title>
      <link>http://example.com/strings</link>
      <description></description>
      <pubDate>Wed, 16 Jan 2013 21:52:35 -0500</pubDate>
    </item>
  </channel>
</rss>`

func TestRSSSorted(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2013-01-16T21:52:35-05:00")
	if err != nil {
		t.Error(err)
	}
	tz := time.FixedZone("EST", -5*60*60)
	now = now.In(tz)

	feed := &Feed{
		Title:       "jmoiron.net blog",
		Link:        &Link{Href: "http://jmoiron.net/blog"},
		Description: "discussion about tech, footie, photos",
		Author:      &Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
		Created:     now,
		Copyright:   "This work is copyright © Benjamin Button",
	}

	feed.Items = []*Item{
		{
			Title:   "Limiting Concurrency in Go",
			Link:    &Link{Href: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
			Created: now.Add(time.Duration(time.Hour * 48)),
		},
		{
			Title:   "Logic-less Template Redux",
			Link:    &Link{Href: "http://jmoiron.net/blog/logicless-template-redux/"},
			Created: now.Add(time.Duration(time.Hour * 24)),
		},
		{
			Title:   "Idiomatic Code Reuse in Go",
			Link:    &Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
			Created: now.Add(time.Duration(time.Hour * 12)),
		},
		{
			Title:   "Never Gonna Give You Up Mp3",
			Link:    &Link{Href: "http://example.com/RickRoll.mp3"},
			Created: now.Add(time.Duration(time.Hour * 10)),
		},
		{
			Title:   "String formatting in Go",
			Link:    &Link{Href: "http://example.com/strings"},
			Created: now,
		}}

	feed.Sort(func(a, b *Item) bool {
		return a.Created.After(b.Created)
	})

	rss, err := feed.ToRss()
	if err != nil {
		t.Errorf("unexpected error encoding RSS: %v", err)
	}

	if rss != rssOutputSorted {
		t.Errorf("Rss not what was expected.  Got:\n%s\n\nExpected:\n%s\n", rss, rssOutputSorted)
	}

	var buf bytes.Buffer
	if err := feed.WriteRss(&buf); err != nil {
		t.Errorf("unexpected error writing RSS: %v", err)
	}
	if got := buf.String(); got != rssOutputSorted {
		t.Errorf("Rss not what was expected.  Got:\n%s\n\nExpected:\n%s\n", got, rssOutputSorted)
	}

}
