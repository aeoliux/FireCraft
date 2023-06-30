package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type MsAuth struct {
	AuthCode     string
	ClientID     string
	AccessToken  string
	RefreshToken string
}

func NewMsAuth(clientID string) (*MsAuth, error) {
	ret := MsAuth{ClientID: clientID}
	if err := ret.GetAuthorizationCode(); err != nil {
		return nil, err
	}
	if err := ret.GetAccessToken(); err != nil {
		return nil, err
	}

	return &ret, nil
}

func (m *MsAuth) GetAuthorizationCode() error {
	httpAddress :=
		"https://login.microsoftonline.com/consumers/oauth2/v2.0/authorize" +
			"?client_id=" + m.ClientID +
			"&response_type=code" +
			"&redirect_uri=http%3A%2F%2Flocalhost%3A8080" +
			"&response_mode=query" +
			"&scope=XboxLive%2Esignin%20offline_access"

	cmd := exec.Command("xdg-open", httpAddress)
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Open '"+httpAddress+"' in your web browser")
	}

	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return err
	}

	conn, err := l.Accept()
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}

	r := `HTTP/1.1 200 OK
Content-Type: text/html

<!DOCTYPE HTML>
<html>
<body>
<script>
window.close()
</script>
</body>
</html>`
	conn.Write([]byte(r))

	conn.Close()
	l.Close()

	response := string(buf)
	code := strings.Split(
		strings.Split(
			strings.Split(response, "\n")[0],
			" ",
		)[1],
		"=",
	)[1]

	if strings.Split(code, "&")[0] == "access_denied" {
		return errors.New("Access denied")
	}

	m.AuthCode = code

	return nil
}

func (m *MsAuth) GetAccessToken() error {
	if m.AuthCode == "" && m.ClientID == "" {
		return errors.New("Missing access token")
	}

	postdata := "client_id=" + m.ClientID +
		"&scope=XboxLive%2Esignin%20offline_access" +
		"&code=" + m.AuthCode +
		"&redirect_uri=http%3A%2F%2Flocalhost%3A8080" +
		"&grant_type=authorization_code"

	requestBody := bytes.NewBufferString(postdata)
	resp, err := http.Post("https://login.microsoftonline.com/consumers/oauth2/v2.0/token", "application/x-www-form-urlencoded", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var j map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(body))
	if err := dec.Decode(&j); err != nil {
		return err
	}

	accessToken, ok := j["access_token"].(string)
	if !ok {
		errorName := j["error"].(string)
		errorCodes := j["error_codes"].([]int)
		return errors.New(fmt.Sprintln(errorName+". Error codes: ", errorCodes))
	}

	m.AccessToken = accessToken
	m.RefreshToken = j["refresh_token"].(string)

	return nil
}
