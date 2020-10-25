package crawler

import (
	"testing"
)

func BenchmarkCrawlUrls(b *testing.B) {
	var result map[string][]string
	urls := []string{
		"https://gobyexample.com/http-servers",
		"https://gobyexample.com/http-servers",
		"https://gobyexample.com/http-servers",
		"https://gobyexample.com/http-servers",
		"https://gobyexample.com/http-servers",
		"https://gobyexample.com/http-servers",
		"https://golang.org/pkg/sync/atomic/",
	}
	for n := 0; n < b.N; n++ {
		result = CrawlUrls(urls)
	}

	b.Log(len(result))
}
