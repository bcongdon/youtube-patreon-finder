package channels

import "testing"

func TestNew(t *testing.T) {
	c, err := New("https://www.youtube.com/feeds/videos.xml?channel_id=UCrYYRuhieWktnpH3UhRahCw", "Foo")
	if err != nil {
		t.Fatalf("New() got err: %v", err)
	}
	wantID := "UCrYYRuhieWktnpH3UhRahCw"
	if c.id != wantID {
		t.Errorf("c.id got %v, want %v", c.id, wantID)
	}

	if c.name != "Foo" {
		t.Errorf("c.name got %v, want %v", c.name, "Foo")
	}
}

func TestRSSURL(t *testing.T) {
	url := "https://www.youtube.com/feeds/videos.xml?channel_id=UCrYYRuhieWktnpH3UhRahCw"
	c, err := New(url, "Foo")
	if err != nil {
		t.Fatalf("New() got err: %v", err)
	}
	if got := c.RSSURL(); got != url {
		t.Errorf("c.RSSURL() got %v, want %v", got, url)
	}
}

func TestNewFailure(t *testing.T) {
	_, err := New("foo", "bar")
	if err == nil {
		t.Errorf(`New("foo", "bar") wanted failure`)
	}
}
