package crawler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/iangregson/spider-go/crawler"
	"github.com/iangregson/spider-go/urls"
)

func TestLocalConcurrency1(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("../fixtures/html")))
	defer server.Close()

	c := crawler.New(1)
	u, _ := url.Parse(server.URL)
	_, u = urls.ValidUrl("/", u)
	c.Crawl(u)

	wanted := 5
	got := c.VisitedCount()

	if got != wanted {
		t.Errorf("Expected %d links visited, got %d", wanted, got)

		for k, v := range c.Visited {
			fmt.Printf("%s %v \n", k, v)
		}
	}
}

func TestLocalConcurrency2(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("../fixtures/html")))
	defer server.Close()

	c := crawler.New(2)
	u, _ := url.Parse(server.URL)
	_, u = urls.ValidUrl("/", u)
	c.Crawl(u)

	wanted := 5
	got := c.VisitedCount()

	if got != wanted {
		t.Errorf("Expected %d links visited, got %d", wanted, got)

		for k, v := range c.Visited {
			fmt.Printf("%s %v \n", k, v)
		}
	}
}

func TestLocalConcurrency16(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("../fixtures/html")))
	defer server.Close()

	c := crawler.New(16)
	u, _ := url.Parse(server.URL)
	_, u = urls.ValidUrl("/", u)
	c.Crawl(u)

	wanted := 5
	got := c.VisitedCount()

	if got != wanted {
		t.Errorf("Expected %d links visited, got %d", wanted, got)

		for k, v := range c.Visited {
			fmt.Printf("%s %v \n", k, v)
		}
	}

	if c.Concurrency() != 16 {
		t.Errorf("Expected %d concurrency, got %d", c.Concurrency(), 16)
	}
}

func TestRemoteConcurrency64(t *testing.T) {
	c := crawler.New(64)
	u, _ := url.Parse("http://crawler-test.com")
	_, u = urls.ValidUrl("/", u)
	c.Crawl(u)

	got := c.VisitedCount()

	if got < 20 {
		t.Errorf("Expected more than 20 links visited, got %d", got)

		for k, v := range c.Visited {
			fmt.Printf("%s %v \n", k, v)
		}
	}
}
