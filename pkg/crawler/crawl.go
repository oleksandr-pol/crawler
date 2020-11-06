/*
	package crawler implements a simple library to concurrently crawl
	HTML pages
*/
package crawler

import (
	"errors"
	"net/http"
	"strings"
	"sync"
)

type Result struct {
	Url   string
	Links []string
	Error error
}

// Crawl fetch HTML pages by endpoints and parse links from fetched page
// It accepts url as parameter, type string.
// It returns urls slice of type string.
// It return error fetched page is not of type HTML
func Crawl(url string) ([]string, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch data")
	}

	notHtml := strings.Index(res.Header.Get("Content-Type"), "text/html") != 0
	if notHtml {
		return nil, errors.New("endpoint does not return html")
	}

	return Links(res.Body), nil
}

// CrawlConcurrent crawl page by url and writes result to a channel.
// It excepts url as parameter, type string.
// It excepts channel as secont parameter, type Result channel.
func CrawlConcurrent(url string, res chan Result) {
	// need to handle error, no blank identifier for errors
	links, err := Crawl(url)
	if err != nil {
		res <- Result{
			Url:   url,
			Error: err,
		}
	} else {
		res <- Result{
			Url:   url,
			Links: links,
		}
	}
}

// CrawlUrls concurrently fetch pages, crawl them and resturns results map.
// It excepts urls slice as first param, type string slice.
// It returns results map. Map key is a url which is crawled, type string.
// Results map values are urls crawled from page. Type string slice.
func CrawlUrls(urls []string) map[string][]string {
	var wg = new(sync.WaitGroup)
	var allResults = make(map[string][]string)
	var resultChan = make(chan Result)

	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			CrawlConcurrent(url, resultChan)
			wg.Done()
		}(url)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		allResults[result.Url] = result.Links
	}

	return allResults
}
