package urls_test

import (
	"net/url"
	"testing"

	"github.com/iangregson/spider-go/urls"
)

func TestValidUrl(t *testing.T) {
	baseUrl, _ := url.Parse("http://localhost:3000")
	u := "https://google.com"
	got, _ := urls.ValidUrl(u, baseUrl)
	want := true

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	u = "https://google.com"
	want = true
	got, _ = urls.ValidUrl(u, baseUrl)

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	u = "monzo.com"
	want = true
	got, _ = urls.ValidUrl(u, baseUrl)

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	u = "http://monzo.com"
	want = true
	got, _ = urls.ValidUrl(u, baseUrl)

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	u = "http://some-subdomain.monzo.com"
	want = true
	got, _ = urls.ValidUrl(u, baseUrl)

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	u = "http://bad subdomain.monzo.com"
	want = false
	got, _ = urls.ValidUrl(u, baseUrl)

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	u = "http://monzo.com/somewhere?odd *&^4 ; Query=123"
	want = true
	got, _ = urls.ValidUrl(u, baseUrl)

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	u = "http://monzo.com/docs/some-page#with-anchor"
	want = true
	got, _ = urls.ValidUrl(u, baseUrl)

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	u = "/docs/some-page#with-anchor"
	want = true
	wantedUrl := "http://localhost:3000/docs/some-page#with-anchor"
	got, validUrl := urls.ValidUrl(u, baseUrl)

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	if validUrl.String() != wantedUrl {
		t.Errorf("got %v, wanted %v, for %s", validUrl.String(), wantedUrl, u)
	}

	u = "#with-anchor"
	want = true
	wantedUrl = "http://localhost:3000/#with-anchor"
	got, validUrl = urls.ValidUrl(u, baseUrl)

	if got != want {
		t.Errorf("got %v, wanted %v, for %s", got, want, u)
	}

	if validUrl.String() != wantedUrl {
		t.Errorf("got %v, wanted %v, for %s", validUrl.String(), wantedUrl, u)
	}
}
