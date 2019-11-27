package ytpf

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/bcongdon/youtube-patreon-finder/pkg/channels"
	"github.com/cheggaaa/pb/v3"
	"github.com/gilliek/go-opml/opml"
)

const parallelism = 5

const mb = 1 << (10 * 3)

type Subscription struct {
	Channel    *channels.Channel `json:"channel"`
	PatreonURL string            `json:"patreon_url"`
	Err        error             `json:"-"`
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

func getPatreonURLForChannel(channelURL string) (string, error) {
	res, err := http.Get(channelURL)
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

func FromOPML(data []byte) ([]*Subscription, error) {
	doc, err := opml.NewOPML(data)
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
	for i := 0; i < parallelism; i++ {
		go func() {
			for c := range inCh {
				patreonLink, err := getPatreonURLForChannel(c.AboutURL())
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

func FromOPMLFile(path string) ([]*Subscription, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return FromOPML(data)
}

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 5*mb))
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading body: %v", err), http.StatusBadRequest)
		return
	}

	subs, err := FromOPML(body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting subscription list: %v", err), http.StatusBadRequest)
		return
	}

	// Filter out subscriptions that we didn't find a Patreon link for.
	filteredSubs := []*Subscription{}
	for _, s := range subs {
		if s.PatreonURL != "" && s.Err == nil {
			filteredSubs = append(filteredSubs, s)
		}
	}

	resp, err := json.Marshal(filteredSubs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
