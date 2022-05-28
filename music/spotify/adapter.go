package spotify

import (
	"partie-bot/requests"

	"golang.org/x/net/html"
)

func URLToText(url string) string {
	response, err := requests.Get(url)
	if err != nil {
		return ""
	}

	tokenizer := html.NewTokenizer(response.Body)
	isTitle := false
	for {
		tt := tokenizer.Next()

		switch {
		case tt == html.ErrorToken:
			return ""

		case tt == html.StartTagToken:
			t := tokenizer.Token()
			isTitle = t.Data == "title"

		case tt == html.TextToken:
			t := tokenizer.Token()

			if isTitle {
				return t.Data
			}
		}
	}
}
