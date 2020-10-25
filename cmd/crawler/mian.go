package main

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type Result struct {
	url   string
	value []byte
}

func main() {

}

func crawl(url string) ([]Result, error) {
	res, err := http.Get("localhost:8000/materials")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch data")
	}

	if res.Header.Get("Content-Type") != "text/html; charset=utf-8" {
		return nil, errors.New("endpoint does not return html")
	}

	urls := getLinks(res.Body)

}

func getLinks(b io.Reader) []string {
	var res []string
	tokenizer := html.NewTokenizer(body)

	for {
		token := tokenizer.Next()

		switch {
		case token == html.ErrorToken:
			return res
		case token == html.StartTagToken:
			t := token.Token()

			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, url := getHref(t)
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
