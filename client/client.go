package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Token      string
}

type Credentials struct {
	AccessToken string `json:"access_token"`
}

func NewClient(baseUrl, clientId, clientSecret string) (*Client, error) {
	values := map[string]string{"grant_type": "client_credentials", "scope": "default"}
	json_data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/oauth/token", baseUrl), bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(clientId, clientSecret)
	req.Header.Add("Accept", "application/json")
	httpClient := http.DefaultClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var credentials Credentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		return nil, err
	}
	return &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		BaseURL:    baseUrl,
		Token:      credentials.AccessToken,
	}, nil
}

func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("Status code: %d, Response body: %s", res.StatusCode, string(body))
	}

	return body, nil
}
