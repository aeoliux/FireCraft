package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type MinecraftAuthentication struct {
	AccessToken    string
	XboxLive       string
	Userhash       string
	XstsToken      string
	MinecraftToken string

	Client *HttpClient
}

func NewMinecraftAuthentication(accessToken string, client *HttpClient) (*MinecraftAuthentication, error) {
	m := MinecraftAuthentication{AccessToken: accessToken, Client: client}
	if err := m.AuthXboxLive(); err != nil {
		return nil, err
	}
	if err := m.AuthXSTS(); err != nil {
		return nil, err
	}
	if err := m.AuthenticateMinecraft(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *MinecraftAuthentication) AuthXboxLive() error {
	reqBody := fmt.Sprintf(`{
	"Properties": {
		"AuthMethod": "RPS",
		"SiteName": "user.auth.xboxlive.com",
		"RpsTicket": "%s",
	},
	"RelyingParty": "http://auth.xboxlive.com",
	"TokenType": "JWT"
}`, m.AccessToken)

	response, err := m.Client.POST("https://user.auth.xboxlive.com/user/authenticate", reqBody, []string{"Content-Type", "Accept"}, []string{"application/json", "application/json"})
	if err != nil {
		return err
	}

	var j map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(response))
	if err := dec.Decode(&j); err != nil {
		return err
	}

	token, ok := j["Token"].(string)
	if !ok {
		return errors.New("Failed to authenticate XBoxLive. Missing token")
	}
	userhash, ok := j["DisplayClaims"].(map[string]interface{})["xui"].([]interface{})[0].(map[string]interface{})["uhs"].(string)
	if !ok {
		return errors.New("Failed to authenticate XBoxLive. Missing userhash")
	}

	m.XboxLive = token
	m.Userhash = userhash

	return nil
}

func (m *MinecraftAuthentication) AuthXSTS() error {
	reqBody := fmt.Sprintf(`{
	"Properties": {
		"SandboxId": "RETAIL",
		"UserTokens": [
			"%s"
		],
	},
	"RelyingParty": "rp://api.minecraftservices.com/",
	"TokenType": "JWT"
}`, m.XboxLive)

	response, err := m.Client.POST("https://xsts.auth.xboxlive.com/xsts/authorize", reqBody, []string{"Content-Type", "Accept"}, []string{"application/json", "application/json"})
	if err != nil {

		return err
	}

	var j map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(response))
	if err := dec.Decode(&j); err != nil {
		return err
	}

	xsts, ok := j["Token"].(string)
	if !ok {
		return errors.New("Failed to authenticate XSTS")
	}

	m.XstsToken = xsts
	return nil
}

func (m *MinecraftAuthentication) AuthenticateMinecraft() error {
	reqBody := fmt.Sprintf(`{"identityToken": "XBL3.0 x=%s;%s","ensureLegacyEnabled": true}`, m.Userhash, m.XstsToken)
	resp, err := m.Client.POST("https://api.minecraftservices.com/authentication/login_with_xbox", reqBody, []string{}, []string{})
	if err != nil {
		return err
	}

	var j map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(resp))
	if err := dec.Decode(&j); err != nil {
		return err
	}

	token, ok := j["access_token"].(string)
	if !ok {
		return errors.New("Failed to authenticate Minecraft. Missing token")
	}

	m.MinecraftToken = token

	return nil
}

// func request(url, body string, headers bool, keys, vals []string) ([]byte, error) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
// 	if err != nil {
// 		return []byte{}, err
// 	}

// 	if headers {
// 		req.Header.Add("Content-Type", "application/json")
// 		req.Header.Add("Accept", "application/json")
// 	}
// 	if len(keys) == len(vals) {
// 		for i := 0; i < len(keys); i++ {
// 			req.Header.Add(keys[i], vals[i])
// 		}
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return []byte{}, err
// 	}
// 	defer resp.Body.Close()

// 	b, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return []byte{}, err
// 	}

// 	return b, nil
// }
