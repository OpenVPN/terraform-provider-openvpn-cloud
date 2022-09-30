package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Host struct {
	Id             string      `json:"id,omitempty"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	InternetAccess string      `json:"internetAccess"`
	SystemSubnets  []string    `json:"systemSubnets"`
	Connectors     []Connector `json:"connectors"`
}

func (c *Client) GetHosts() ([]Host, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/hosts", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var hosts []Host
	err = json.Unmarshal(body, &hosts)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func (c *Client) GetHostByName(name string) (*Host, error) {
	hosts, err := c.GetHosts()
	if err != nil {
		return nil, err
	}
	for _, h := range hosts {
		if h.Name == name {
			return &h, nil
		}
	}
	return nil, nil
}

func (c *Client) GetHostById(hostId string) (*Host, error) {
	hosts, err := c.GetHosts()
	if err != nil {
		return nil, err
	}
	for _, h := range hosts {
		if h.Id == hostId {
			return &h, nil
		}
	}
	return nil, nil
}

func (c *Client) CreateHost(host Host) (*Host, error) {
	hostJson, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/hosts", c.BaseURL), bytes.NewBuffer(hostJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var h Host
	err = json.Unmarshal(body, &h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (c *Client) UpdateHost(host Host) error {
	hostJson, err := json.Marshal(host)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/beta/hosts/%s", c.BaseURL, host.Id), bytes.NewBuffer(hostJson))
	if err != nil {
		return err
	}
	_, err = c.DoRequest(req)
	return err
}

func (c *Client) DeleteHost(hostId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/hosts/%s", c.BaseURL, hostId), nil)
	if err != nil {
		return err
	}
	_, err = c.DoRequest(req)
	return err
}
