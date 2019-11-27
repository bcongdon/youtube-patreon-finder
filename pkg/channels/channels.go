package channels

import (
	"fmt"
	"regexp"
)

var channelIdRegex = regexp.MustCompile("channel_id=(.*)$")

type Channel struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (c *Channel) RSSURL() string {
	return fmt.Sprintf(
		"https://www.youtube.com/feeds/videos.xml?channel_id=%s", c.ID)
}

func (c *Channel) AboutURL() string {
	return fmt.Sprintf(
		"https://www.youtube.com/channel/%s/about", c.ID)
}

func New(url, name string) (*Channel, error) {
	matches := channelIdRegex.FindStringSubmatch(url)

	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid channel url: %q", url)
	}
	return &Channel{
		ID:   matches[1],
		Name: name,
	}, nil
}
