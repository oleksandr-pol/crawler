package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/oleksandr-pol/crawler/pkg/crawler"
)

func main() {
	done := make(chan os.Signal)
	signal.Notify(done, os.Interrupt)

	defer close(done)

	links, err := crawler.Crawl("http://localhost:8000/materials")
	if err != nil {
		fmt.Print(err.Error())
	}

	// first version
	// result := crawler.CrawlUrls(links)

	result := crawler.CrawlPipe(done, links)

	for key, val := range result {
		fmt.Println(key, len(val))
	}
}
