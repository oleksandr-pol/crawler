package crawler

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Links parses links from reader.
// It accepts io.Reader as argument.
// It return slice of links, type slice of string.
func Links(b io.Reader) []string { // getters does not need get as prefix
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

			ok, url := Href(token)
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
