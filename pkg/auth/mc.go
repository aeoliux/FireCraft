package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MinecraftAuthentication struct {
	AccessToken    string
	XboxLive       string
	Userhash       string
	XstsToken      string
	MinecraftToken string
}

func NewMinecraftAuthentication(accessToken string) (*MinecraftAuthentication, error) {
	m := MinecraftAuthentication{AccessToken: accessToken}
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
		"RpsTicket": "d=%s",
	},
	"RelyingParty": "http://auth.xboxlive.com",
	"TokenType": "JWT"
}`, m.AccessToken)

	response, err := request("https://user.auth.xboxlive.com/user/authenticate", reqBody, true)
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

	response, err := request("https://xsts.auth.xboxlive.com/xsts/authorize", reqBody, true)
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
	reqBody := fmt.Sprintf(`{"identityToken": "XBL3.0 x=%s;%s"}`, m.Userhash, m.XstsToken)
	resp, err := request("https://api.minecraftservices.com/authentication/login_with_xbox", reqBody, false)
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

func request(url, body string, headers bool) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	if err != nil {
		return []byte{}, err
	}

	if headers {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return b, nil
}
