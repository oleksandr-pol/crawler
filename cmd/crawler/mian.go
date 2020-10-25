package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	res, err := crawl("http://localhost:8000/materials")
	if err != nil {
		fmt.Print(err.Error())
	}

	for _, url := range res {
		urls, err := crawl(url)
		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Println(urls)
	}
}

func crawl(url string) ([]string, error) {
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

	return getLinks(res.Body), nil
}

func getLinks(b io.Reader) []string {
	var res []string
	tokenizer := html.NewTokenizer(b)

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			return res
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()

			isAnchor := token.Data == "a"
			if !isAnchor {
				continue
			}

			ok, url := getHref(token)
			if !ok {
				continue
			}

			hasProtocol := strings.Index(url, "http") == 0
			if hasProtocol {
				res = append(res, url)
			}
		}
	}
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	return
}
