package auth

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Jar struct {
	jar map[string][]*http.Cookie
}

func (j *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.jar[u.Host] = cookies
}

func (j *Jar) Cookies(u *url.URL) []*http.Cookie {
	return j.jar[u.Host]
}

type HttpClient struct {
	client       *http.Client
	req          *http.Request
	LastRedirect string
}

func NewHttpClient() *HttpClient {
	ht := &HttpClient{}
	ht.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			ht.LastRedirect = req.URL.String()
			return nil
		},
	}
	jar := &Jar{}
	jar.jar = make(map[string][]*http.Cookie)
	ht.client.Jar = jar
	return ht
}

func (h *HttpClient) GET(url string, keys, vals []string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if len(keys) == len(vals) {
		for i := range keys {
			req.Header.Add(keys[i], vals[i])
		}
	}

	return h.request(req)
}

func (h *HttpClient) POST(url, body string, keys, vals []string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	if len(keys) == len(vals) {
		for i := range keys {
			req.Header.Add(keys[i], vals[i])
		}
	}

	return h.request(req)
}

func (h *HttpClient) request(req *http.Request) ([]byte, error) {
	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return bodyBytes, nil
}
