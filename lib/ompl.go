package lib

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bcongdon/youtube-patreon-finder/lib/channels"
	"github.com/gilliek/go-opml/opml"
)

type Subscription struct {
	Channel    *channels.Channel
	PatreonURL string
}

func parsePatreonLinkFromRedirect(redirect string) (string, bool) {
	u, err := url.Parse(redirect)
	if err != nil {
		return "", false
	}
	q := u.Query().Get("q")
	if strings.Contains(q, "patreon.com") {
		return q, true
	}
	return "", false
}

func GetPatreonURL(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", nil
	}
	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	var patreonLink string
	doc.Find(".channel-links-item").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		href, ok := s.Find("a").Attr("href")
		if !ok {
			return true
		}
		patreonLink, ok = parsePatreonLinkFromRedirect(href)
		return !ok
	})

	return patreonLink, nil
}

func FromFile(path string) ([]*Subscription, error) {
	doc, err := opml.NewOPMLFromFile(path)
	if err != nil {
		log.Fatal(err)
	}

	rootOutlines := doc.Outlines()
	subscriptions := rootOutlines[0].Outlines

	var out []*Subscription
	for _, sub := range subscriptions {
		c, err := channels.New(sub.XMLURL, sub.Title)
		if err != nil {
			return nil, err
		}
		patreonLink, err := GetPatreonURL(c.AboutURL())
		fmt.Println(c, patreonLink)
		if err != nil {
			fmt.Println(err)
			continue
		}
		out = append(out, &Subscription{c, patreonLink})
	}

	return out, nil
}
