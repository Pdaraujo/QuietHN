package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiBase = "https://hacker-news.firebaseio.com/v0"
	topStories = "/topstories.json"
)

// Item represents a single item returned by the HN API. This can have a type
// of "story", "comment", or "job" (and probably more values), and one of the
// URL or Text fields will be set, but not both.
//
// For the purpose of this exercise, we only care about items where the
// type is "story", and the URL is set.
type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	// Only one of these should exist
	Text string `json:"text"`
	URL  string `json:"url"`
}

type Client struct {
	apiBase string
	topStories string
	topItemsURL string
}

// Making the Client zero value useful without forcing users to do something
// like `NewClient()`
func (c *Client) setDefaults() {
	if c.apiBase == "" {
		c.apiBase = apiBase
	}

	if c.topStories == "" {
		c.topStories = topStories
	}

	c.topItemsURL = fmt.Sprintf("%s%s", c.apiBase, c.topStories)
}

func (c *Client) createItemURL(id int) string {
	c.setDefaults()
	return fmt.Sprintf("%s/item/%d.json", c.apiBase, id)
}

// TopItems returns the ids of roughly 450 top items in decreasing order. These
// should map directly to the top 450 things you would see on HN if you visited
// their site and kept going to the next page.
//
// TopItems does not filter out job listings or anything else, as the type of
// each item is unknown without further API calls.
func (c *Client) TopItems() ([]int, error) {
	c.setDefaults()
	resp, err := http.Get(c.topItemsURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var ids []int
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// GetItem will return the Item defined by the provided ID.
func (c *Client) GetItem(id int) (Item, error) {
	var item Item
	resp, err := http.Get(c.createItemURL(id))
	if err != nil {
		return item, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)
	if err != nil {
		return item, err
	}
	return item, nil
}