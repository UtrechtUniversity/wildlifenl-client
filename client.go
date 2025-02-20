package wildlifenl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func New(baseURL string) *Client {
	c := new(Client)
	c.baseURL = baseURL
	c.webclient = &http.Client{}
	return c
}

type Client struct {
	baseURL    string
	webclient  *http.Client
	credential *Credential
}

func (c *Client) Authenticate(appName string, email string) error {
	body := make(map[string]any)
	body["displayNameApp"] = appName
	body["email"] = email
	payload, _ := json.Marshal(body)
	if _, err := c.Call(http.MethodPost, "/auth/", payload, nil); err != nil {
		return fmt.Errorf("cannot authenticate: %w", err)
	}
	return nil
}

func (c *Client) Authorize(email string, code string) (*Credential, error) {
	body := make(map[string]any)
	body["email"] = email
	body["code"] = code
	payload, _ := json.Marshal(body)
	data, err := c.Call(http.MethodPut, "/auth/", payload, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot authorize: %w", err)
	}
	credential := new(Credential)
	if err := json.Unmarshal(data, credential); err != nil {
		return nil, fmt.Errorf("cannot parse credential from authorize response: %w", err)
	}
	return credential, nil
}

func (c *Client) Call(method string, path string, body []byte, credential *Credential) ([]byte, error) {
	endpoint, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return nil, fmt.Errorf("cannot join baseURL and path: %w", err)
	}
	request, err := http.NewRequest(method, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("cannot prepare http request: %w", err)
	}
	if credential != nil {
		request.Header.Add("Authorization", "Bearer "+c.credential.Token)
	}
	response, err := c.webclient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("cannot process http request: %w", err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read http response: %w", err)
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("http resonse has status %v: %w", response.Status, errors.New(string(data)))
	}
	return data, nil
}
