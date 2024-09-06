package requests

import (
	"net/http"
	"os"
)

// func SpotifyGet(endpoint string, queryParams map[string]string) (*http.Response, error) {
// 	u, err := url.Parse(endpoint)
// 	if err != nil {
// 		return nil, err
// 	}

// 	q := u.Query()
// 	for key, value := range queryParams {
// 		q.Set(key, value)
// 	}

// 	u.RawQuery = q.Encode()

// 	response, err := http.Get(u.String())
// 	return response, err
// }

func SpotifyGet(endpoint string, queryParams map[string]string) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", endpoint, nil)
	request.Header.Set("Authorization", "Bearer "+os.Getenv("SPOTIFY_TOKEN"))

	q := request.URL.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}

	request.URL.RawQuery = q.Encode()

	response, err := client.Do(request)
	return response, err
}
