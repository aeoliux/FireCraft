package auth

import (
	"errors"
	"fmt"
	"strings"
)

type Authentication struct {
	HtClient *HttpClient

	MsAccessToken  string
	MsRefreshToken *string
}

func NewAuthentication(email, password string) (*Authentication, error) {
	au := Authentication{}
	au.HtClient = NewHttpClient()

	bodystr, err := au.HtClient.GET(
		"https://login.live.com/oauth20_authorize.srf?client_id=000000004C12AE6F&redirect_uri=https://login.live.com/oauth20_desktop.srf&scope=service::user.auth.xboxlive.com::MBI_SSL&display=touch&response_type=token&locale=en",
		[]string{},
		[]string{},
	)
	if err != nil {
		return nil, err
	}

	fftagVal := extractsFTTagValue(string(bodystr))
	urlPost := extractUrlPost(string(bodystr))

	post := fmt.Sprintf("login=%s&loginfmt=%s&passwd=%s&PPFT=%s", email, email, password, *fftagVal)
	resp, err := au.HtClient.POST(*urlPost, post, []string{"Content-Type"}, []string{"application/x-www-form-urlencoded"})
	if err != nil {
		return nil, err
	}

	if au.HtClient.LastRedirect != "" {
		ma := UrlEncodedToMap(au.HtClient.LastRedirect)
		accessToken, ok := ma["access_token"]
		if !ok {
			return nil, errors.New("Login failed: " + string(resp))
		}
		au.MsAccessToken = accessToken

		refreshToken, ok := ma["refresh_token"]
		if ok {
			au.MsRefreshToken = &refreshToken
		}
	} else {
		return nil, errors.New("Login failed: " + string(resp))
	}

	return &au, nil
}

func extractsFTTagValue(str string) *string {
	for i := 0; i < len(str); i++ {
		if strings.HasPrefix(str[i:], "sFTTag:'<") {
			split1 := strings.Split(str[i:], "'")
			trim := strings.Trim(split1[1], "'")

			split2 := strings.Split(trim, " ")
			for _, j := range split2 {
				if strings.HasPrefix(j, "value=\"") {
					split3 := strings.Split(j, "=")
					ret := strings.TrimSuffix(strings.TrimPrefix(split3[1], "\""), "\"/>")
					return &ret
				}
			}
		}
	}

	return nil
}

func extractUrlPost(str string) *string {
	for i := 0; i < len(str); i++ {
		if strings.HasPrefix(str[i:], "urlPost:'") {
			split := strings.Split(str[i:], "'")
			return &split[1]
		}
	}

	return nil
}
