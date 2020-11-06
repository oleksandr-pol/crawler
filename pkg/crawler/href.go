package crawler

import "golang.org/x/net/html"

// Href gets href attribute from html token.
// It excepts html token as parameter.
// It return flag tat indicates if html token has href atts.
// It returns href attribute value.
func Href(t html.Token) (ok bool, href string) { // getters does not need get as prefix
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	return
}
