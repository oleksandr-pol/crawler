package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oleksandr-pol/crawler/pkg/crawler"
)

func main() {
	done := make(chan struct{})
	stopCrawl(done)

	links, err := crawler.Crawl("http://localhost:8000/materials")
	if err != nil {
		fmt.Print(err.Error())
	}

	result := crawler.CrawlPipe(done, links)

	for key, val := range result {
		fmt.Println(key, len(val))
	}
}

func stopCrawl(done chan struct{}) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		close(done)
		os.Exit(0)
		return
	}()
}
