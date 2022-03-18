package requests

import (
	"fmt"
	"net/http"
	"net/url"
)

func YoutubeSearch(query string) (*http.Response, error) {
	searchUrl := url.URL{
		Host:     "localhost:3000",
		Scheme:   "http",
		Path:     "/",
		RawQuery: "q=" + url.QueryEscape(query),
	}

	fmt.Println("Searching youtube for: ", searchUrl.String())

	return http.Get(searchUrl.String())
}
