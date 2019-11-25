package lib

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/bcongdon/youtube-patreon-finder/lib/channels"
	"github.com/cheggaaa/pb/v3"
	"github.com/gilliek/go-opml/opml"
)

type Subscription struct {
	Channel    *channels.Channel
	PatreonURL string
	Err        error
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
		return nil, err
	}

	rootOutlines := doc.Outlines()
	subscriptions := rootOutlines[0].Outlines
	bar := pb.StartNew(len(subscriptions))

	var wg sync.WaitGroup
	inCh := make(chan *channels.Channel)
	outCh := make(chan *Subscription, 10)

	// Patreon link fetchers
	for i := 0; i < 5; i++ {
		go func() {
			for c := range inCh {
				patreonLink, err := GetPatreonURL(c.AboutURL())
				outCh <- &Subscription{c, patreonLink, err}
			}
		}()
	}

	go func() {
		defer close(inCh)
		defer close(outCh)
		for _, sub := range subscriptions {
			c, err := channels.New(sub.XMLURL, sub.Title)
			if err != nil {
				continue
			}
			wg.Add(1)
			inCh <- c
		}
		wg.Wait()
	}()

	var subs []*Subscription
	for c := range outCh {
		subs = append(subs, c)
		wg.Done()
		bar.Increment()
	}
	bar.Finish()

	return subs, nil
}
