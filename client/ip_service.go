package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Range struct {
	LowerValue int `json:"lowerValue"`
	UpperValue int `json:"upperValue"`
}

type CustomIpServiceType struct {
	IcmpType []Range `json:"icmpType"`
	Port     []Range `json:"port"`
	Protocol string  `json:"protocol"`
}

type IpServiceConfig struct {
	CustomServiceTypes []*CustomIpServiceType `json:"customServiceTypes"`
	ServiceTypes       []string               `json:"serviceTypes"`
}

type IpService struct {
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	NetworkItemType string           `json:"networkItemType"`
	NetworkItemId   string           `json:"networkItemId"`
	Id              string           `json:"id"`
	Type            string           `json:"type"`
	Routes          []*Route         `json:"routes"`
	Config          *IpServiceConfig `json:"config"`
}

type IpServicePageResponse struct {
	Content          []IpService `json:"content"`
	NumberOfElements int         `json:"numberOfElements"`
	Page             int         `json:"page"`
	Size             int         `json:"size"`
	Success          bool        `json:"success"`
	TotalElements    int         `json:"totalElements"`
	TotalPages       int         `json:"totalPages"`
}

func (c *Client) GetIpServicesByPage(page int, pageSize int) (IpServicePageResponse, error) {
	endpoint := fmt.Sprintf("%s/api/beta/ip-services/page?page=%d&size=%d", c.BaseURL, page, pageSize)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return IpServicePageResponse{}, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return IpServicePageResponse{}, err
	}

	var response IpServicePageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return IpServicePageResponse{}, err
	}
	return response, nil
}

func (c *Client) GetAllIpServices() ([]IpService, error) {
	var allIpServices []IpService
	page := 1
	pageSize := 10

	for {
		response, err := c.GetIpServicesByPage(page, pageSize)
		if err != nil {
			return nil, err
		}

		allIpServices = append(allIpServices, response.Content...)
		if page >= response.TotalPages {
			break
		}
		page++
	}
	return allIpServices, nil
}

func (c *Client) CreateIpService(ipService *IpService) (*IpService, error) {
	ipServiceJson, err := json.Marshal(ipService)
	if err != nil {
		return nil, err
	}

	params := networkUrlParams(ipService.NetworkItemType, ipService.NetworkItemId)
	endpoint := fmt.Sprintf("%s/api/beta/ip-services?%s", c.BaseURL, params.Encode())

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(ipServiceJson))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var s IpService
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) UpdateIpService(id string, service *IpService) (*IpService, error) {
	serviceJson, err := json.Marshal(service)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/api/beta/ip-services/%s?%s", c.BaseURL, id)

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(serviceJson))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var s IpService
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) DeleteIpService(ipServiceId string) error {
	endpoint := fmt.Sprintf("%s/api/beta/ip-services/%s?%s", c.BaseURL, ipServiceId)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func networkUrlParams(networkItemType string, networkItemId string) url.Values {
	params := url.Values{}
	params.Add("networkItemId", networkItemId)
	params.Add("networkItemType", networkItemType)
	return params
}
