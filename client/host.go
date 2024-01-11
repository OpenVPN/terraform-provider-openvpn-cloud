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

type HostPageResponse struct {
	Content          []Host `json:"content"`
	NumberOfElements int    `json:"numberOfElements"`
	Page             int    `json:"page"`
	Size             int    `json:"size"`
	Success          bool   `json:"success"`
	TotalElements    int    `json:"totalElements"`
	TotalPages       int    `json:"totalPages"`
}

func (c *Client) GetHostsByPage(page int, size int) (HostPageResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/hosts/page?page=%d&size=%d", c.BaseURL, page, size), nil)
	if err != nil {
		return HostPageResponse{}, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return HostPageResponse{}, err
	}

	var response HostPageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return HostPageResponse{}, err
	}
	return response, nil
}

func (c *Client) GetAllHosts() ([]Host, error) {
	var allHosts []Host
	pageSize := 10
	page := 1

	for {
		response, err := c.GetHostsByPage(page, pageSize)
		if err != nil {
			return nil, err
		}

		allHosts = append(allHosts, response.Content...)

		if page >= response.TotalPages {
			break
		}
		page++
	}
	return allHosts, nil
}

func (c *Client) GetHostByName(name string) (*Host, error) {
	hosts, err := c.GetAllHosts()
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
	hosts, err := c.GetAllHosts()
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
