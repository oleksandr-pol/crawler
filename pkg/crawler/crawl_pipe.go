package crawler

import (
	"sync"
)

func genUrlsChanStage(done <-chan struct{}, urls []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, url := range urls {
			select {
			case out <- url:
			case <-done:
				return
			}
		}
		close(out)
	}()

	return out
}

func crawlStage(done <-chan struct{}, in <-chan string) <-chan Result {
	var wg sync.WaitGroup
	out := make(chan Result)

	for url := range in {
		go func(curUrl string) {
			wg.Add(1)
			var res Result
			links, err := Crawl(curUrl)
			if err != nil {
				res = Result{
					Url:   curUrl,
					Error: err,
				}
			}

			res = Result{
				Url:   curUrl,
				Links: links,
			}

			select {
			case out <- res:
			case <-done:
				return
			}
			wg.Done()
		}(url)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func resultStage(in <-chan Result) map[string][]string {
	var allResults = make(map[string][]string)

	for result := range in {
		allResults[result.Url] = result.Links
	}

	return allResults
}

func CrawlPipe(done <-chan struct{}, urls []string) map[string][]string {
	urlChan := genUrlsChanStage(done, urls)
	outChan := crawlStage(done, urlChan)
	result := resultStage(outChan)

	return result
}
