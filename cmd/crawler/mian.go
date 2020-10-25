package main

import (
	"fmt"

	"github.com/oleksandr-pol/crawler/pkg/crawler"
)

func main() {
	links, err := crawler.Crawl("http://localhost:8000/materials")
	if err != nil {
		fmt.Print(err.Error())
	}

	result := crawler.CrawlUrls(links)

	for key, val := range result {
		fmt.Println(key, len(val))
	}
}
