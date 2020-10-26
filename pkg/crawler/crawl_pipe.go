package crawler

import "os"

func genUrlsChanStage(done <-chan os.Signal, urls []string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, url := range urls {

			select {
			case out <- url:
			case <-done:
				return
			}
		}
	}()
	return out
}

func crawlStage(done <-chan os.Signal, in <-chan string) <-chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		for url := range in {
			links, err := Crawl(url)

			if err != nil {
				return
			}

			select {
			case out <- Result{
				Url:   url,
				Links: links,
			}:
			case <-done:
				return
			}
		}
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

func CrawlPipe(done <-chan os.Signal, urls []string) map[string][]string {
	urlChan := genUrlsChanStage(done, urls)
	outChan := crawlStage(done, urlChan)
	result := resultStage(outChan)

	return result
}
