package crawler

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/semaphore"

	"github.com/iangregson/spider-go/parser"
	"github.com/iangregson/spider-go/urls"
)

func (c *Crawler) Crawl2(baseUrl *url.URL) {
	frontier := make([]*url.URL, 1, c.maxUrlFrontierSize)

	// enqueue seed url
	frontier[0] = baseUrl

	// sync
	var mx sync.Mutex
	var wg sync.WaitGroup
	ctx := context.TODO()
	pool := semaphore.NewWeighted(int64(c.concurrency))

	// while we've got a url frontier, do the work
	mx.Lock()
	for len(frontier) > 0 {
		mx.Unlock()

		if err := pool.Acquire(ctx, 1); err != nil {
			log.Error().Err(err).Msg("Failed to acquire semaphore")
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer pool.Release(1)

			mx.Lock()
			u := frontier[0]
			frontier = frontier[1:]
			if u == nil {
				mx.Unlock()
				return
			}
			if c.seen[u.String()] {
				mx.Unlock()
				return
			}
			c.seen[u.String()] = true
			mx.Unlock()

			links, err := fetchLinks(u)
			if err != nil {
				return
			}

			mx.Lock()
			c.Visited[u.String()] = true
			for _, link := range links {
				log.Printf("queueing %v", link)
				frontier = append(frontier, link)
			}
			mx.Unlock()
		}()

		wg.Wait()
		mx.Lock()
	}
}

func fetchLinks(u *url.URL) ([]*url.URL, error) {
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpClient.Get(u.String())
	if err != nil {
		log.Error().Err(err).Msg("couldn't fetch url")
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := errors.New("http response not ok")
		log.Error().Err(err)
		return nil, err
	}

	log.Info().Msgf("Visited %s", u.String())

	links := parser.ParseLinks(resp.Body, u)
	linksOnSameHost := []*url.URL{}
	for _, link := range links {
		if urls.SameHost(u, link) {
			linksOnSameHost = append(linksOnSameHost, link)
		}
	}

	return linksOnSameHost, nil
}
