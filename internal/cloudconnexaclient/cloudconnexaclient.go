package cloudconnexaclient

import (
	"context"
	"github.com/openvpn/cloudconnexa-go-client/v2/cloudconnexa"
)

type Config struct {
	UserAgent    string
	BaseURL      string
	ClientID     string
	ClientSecret string
}

func (c *Config) Client(ctx context.Context) (*cloudconnexa.Client, error) {
	var cl *cloudconnexa.Client
	var err error
	if err != nil {
		return nil, err
	}
	cl, err = cloudconnexa.NewClient("", "", "")
	if err != nil {
		return nil, err
	}

	cl.UserAgent = c.UserAgent

	return cl, nil
}
