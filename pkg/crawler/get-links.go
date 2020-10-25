package crawler

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func GetLinks(b io.Reader) []string {
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

			ok, url := GetHref(token)
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
