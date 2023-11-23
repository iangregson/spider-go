package crawler

import (
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/iangregson/spider-go/parser"
	"github.com/iangregson/spider-go/urls"
)

type Crawler struct {
	seen               map[string]bool
	Visited            map[string]bool
	concurrency        int
	urlFrontierSize    int
	maxUrlFrontierSize int
}

type crawlResult struct {
	visitedUrl *url.URL
	visitedOk  bool
	linkedUrls []*url.URL
}

func New(concurrency int) *Crawler {
	var c = &Crawler{
		seen:               make(map[string]bool),
		Visited:            make(map[string]bool),
		concurrency:        concurrency,
		maxUrlFrontierSize: 1e6,
		urlFrontierSize:    0,
	}

	return c
}

func (c *Crawler) Crawl(baseUrl *url.URL) {
	frontier := make(chan *url.URL, c.maxUrlFrontierSize)
	crawlResults := make(chan crawlResult, c.concurrency)

	// start worker pool
	for i := 1; i <= c.concurrency; i++ {
		go func() {
			for {
				select {
				case u := <-frontier:
					if u != nil {
						visit(u, crawlResults)
					}
				}
			}
		}()
	}

	// enqueue seed url
	c.enqueue([]*url.URL{baseUrl}, frontier)

	// while we've got a url frontier, do the work
	for c.urlFrontierSize > 0 {
		select {
		case crawlResult := <-crawlResults:
			c.urlFrontierSize--
			c.seen[crawlResult.visitedUrl.String()] = true
			if crawlResult.visitedOk {
				c.Visited[crawlResult.visitedUrl.String()] = true
			}
			c.enqueue(crawlResult.linkedUrls, frontier)
		}
	}

	// when the url frontier is met, we're done
	close(frontier)
	close(crawlResults)
}

func (c *Crawler) enqueue(urls []*url.URL, frontier chan *url.URL) {
	for _, u := range urls {
		if u == nil || c.seen[u.String()] {
			continue
		}

		select {
		case frontier <- u:
			c.seen[u.String()] = true
			c.urlFrontierSize++
		default:
			// url dropped
			log.Info().Msg("url dropped from frontier because frontier full.")
		}
	}
}

func (c *Crawler) VisitedCount() int {
	return len(c.Visited)
}

func (c *Crawler) Concurrency() int {
	return c.concurrency
}
func visit(u *url.URL, crawlResults chan crawlResult) {

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpClient.Get(u.String())
	if err != nil {
		log.Error().Err(err).Msg("couldn't fetch url")
		crawlResults <- crawlResult{
			visitedUrl: u,
			visitedOk:  false,
			linkedUrls: []*url.URL{},
		}
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Msg("http response not ok")
		crawlResults <- crawlResult{
			visitedUrl: u,
			visitedOk:  false,
			linkedUrls: []*url.URL{},
		}
		return
	}

	log.Info().Msgf("Visited %s", u.String())

	links := parser.ParseLinks(resp.Body, u)
	var linksOnSameHost = []*url.URL{}
	for _, link := range links {
		if urls.SameHost(u, link) {
			linksOnSameHost = append(linksOnSameHost, link)
		}
	}

	crawlResults <- crawlResult{
		visitedUrl: u,
		visitedOk:  true,
		linkedUrls: linksOnSameHost,
	}
}
