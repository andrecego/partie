package requests

import (
	"net/http"
)

func Get(url string) (*http.Response, error) {
	response, err := http.Get(url)
	return response, err
}
