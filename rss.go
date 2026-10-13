package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create an HTTP GET request for the RSS feed.
	// The context allows the request to be cancelled or timed out by the caller.
	request, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	// Identify our application to the web server making the request.
	request.Header.Set("User-Agent", "gator")

	// Create an HTTP client and use it to send the request.
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	// Ensure the response body is closed when this function finishes.
	// This releases the network resources associated with the request.
	defer response.Body.Close()

	// Read the entire HTTP response body into memory as a slice of bytes.
	// At this point, data contains the raw RSS XML returned by the website.
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Create an empty RSSFeed struct that will hold the parsed XML data.
	feed := RSSFeed{}

	// Parse the XML bytes and populate the RSSFeed struct.
	// The xml struct tags on RSSFeed and RSSItem determine where
	// each XML element is stored.
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	// Convert HTML entities in the channel title and description
	// into their normal characters.
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Process every article/item in the feed and decode HTML entities
	// in each item's title and description.
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	// Return a pointer to the completed RSSFeed.
	// nil indicates that no error occurred.
	return &feed, nil
}
