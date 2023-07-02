package auth

import "strings"

func UrlEncodedToMap(url string) map[string]string {
	m := make(map[string]string)
	if strings.Contains(url, "?") {
		for i := range url {
			if url[i] == '?' {
				url = url[i+1:]
				break
			}
		}
	}
	if strings.Contains(url, "#") {
		for i := range url {
			if url[i] == '#' {
				url = url[i+1:]
				break
			}
		}
	}

	split := strings.Split(url, "&")
	for _, j := range split {
		split2 := strings.Split(j, "=")
		m[split2[0]] = split2[1]
	}

	return m
}
