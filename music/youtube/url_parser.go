package youtube

import "regexp"

func MatchURL(query string) bool {
	return regexp.
		MustCompile(`^(https?:\/\/)?(www\.)?(youtube\.com|youtu\.be)\/.+$`).
		Match([]byte(query))
}
