package auth

import "errors"

type Authentication struct {
	HtClient *HttpClient

	MsAccessToken  string
	MsRefreshToken *string
}

func NewAuthentication(accessLink string) (*Authentication, error) {
	au := Authentication{}
	au.HtClient = NewHttpClient()

	m := UrlEncodedToMap(accessLink)
	acc, ok := m["access_token"]
	if !ok {
		return nil, errors.New("missing access token")
	}

	au.MsAccessToken = acc

	return &au, nil
}
