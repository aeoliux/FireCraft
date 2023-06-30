package downloader

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(url, output string) ([]byte, error) {
	img, err := FetchFile(url)
	if err != nil {
		return []byte{}, err
	}

	return img, os.WriteFile(output, img, 0644)
}

func FetchFile(url string) ([]byte, error) {
	req, err := http.Get(url)
	if err != nil {
		return []byte{}, nil
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return []byte{}, nil
	}

	return body, nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return true
}
