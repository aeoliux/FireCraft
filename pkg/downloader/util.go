package downloader

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
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

func CheckSumByPath(path, sha1check string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	return CheckSum(data, sha1check)
}

func CheckSum(data []byte, sha1check string) bool {
	sha1Engine := sha1.New()
	sha1Engine.Write(data)
	sha1Hash := hex.EncodeToString(sha1Engine.Sum(nil))

	return sha1Hash == sha1check
}

func DownloadAndCheck(out, url, sha1check string) error {
	outDir := path.Dir(out)
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return err
	}

	data := []byte{}
	var err error
	for i := 0; i < 3; i++ {
		data, err = DownloadFile(url, out)
		if err != nil {
			return err
		}

		if CheckSum(data, sha1check) {
			return nil
		}
	}

	return errors.New("checksum mismatch, downloaded 3 times from " + url)
}
