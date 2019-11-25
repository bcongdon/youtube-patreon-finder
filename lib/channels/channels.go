package channels

import (
	"fmt"
	"regexp"
)

var channelIdRegex = regexp.MustCompile("channel_id=(.*)$")

type Channel struct {
	name string
	id   string
}

func (c *Channel) RSSURL() string {
	return fmt.Sprintf(
		"https://www.youtube.com/feeds/videos.xml?channel_id=%s", c.id)
}

func (c *Channel) AboutURL() string {
	return fmt.Sprintf(
		"https://www.youtube.com/channel/%s/about", c.id)
}

func (c *Channel) Name() string {
	return c.name
}

func New(url, name string) (*Channel, error) {
	matches := channelIdRegex.FindStringSubmatch(url)

	if len(matches) < 2 {
		return nil, fmt.Errorf("invalid channel url: %q", url)
	}
	return &Channel{
		id:   matches[1],
		name: name,
	}, nil
}
