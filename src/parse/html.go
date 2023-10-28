package parse

import (
	"io"

	"golang.org/x/net/html"
)

func PageLinks(reader io.Reader) (links []string, err error) {
	tokenizer := html.NewTokenizer(reader)
	if err != nil {
		return
	}
	for {
		if tt := tokenizer.Next(); tt == html.ErrorToken {
			return
		}
		token := tokenizer.Token()
		if token.Data == "a" && (token.Type == html.StartTagToken || token.Type == html.SelfClosingTagToken) {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
	}
}
