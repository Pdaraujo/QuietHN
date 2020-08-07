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

func (c *Client) getItemURL(id int) string {
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