package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	HTTPClient  *http.Client
	BaseURL     string
	Token       string
	RateLimiter *rate.Limiter
}

type Credentials struct {
	AccessToken string `json:"access_token"`
}

func NewClient(baseUrl, clientId, clientSecret string) (*Client, error) {
	if clientId == "" || clientSecret == "" {
		return nil, ErrCredentialsRequired
	}

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
	body, err := io.ReadAll(resp.Body)
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
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		BaseURL:     baseUrl,
		Token:       credentials.AccessToken,
		RateLimiter: rate.NewLimiter(rate.Every(1*time.Second), 5),
	}, nil
}

func (c *Client) DoRequest(req *http.Request) ([]byte, error) {
	err := c.RateLimiter.Wait(context.Background())
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("Status code: %d, Response body: %s", res.StatusCode, string(body))
	}

	return body, nil
}
