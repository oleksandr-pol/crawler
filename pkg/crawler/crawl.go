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

	return GetLinks(res.Body), nil
}

func CrawlConcurrent(url string, res chan Result) {
	links, _ := Crawl(url)
	res <- Result{
		Url:   url,
		Links: links,
	}
}

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
