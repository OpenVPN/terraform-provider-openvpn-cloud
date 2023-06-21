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

type CustomServiceType struct {
	IcmpType []Range `json:"icmpType"`
	Port     []Range `json:"port"`
	Protocol string  `json:"protocol"`
}

type ServiceConfig struct {
	CustomServiceTypes []*CustomServiceType `json:"customServiceTypes"`
	ServiceTypes       []string             `json:"serviceTypes"`
}

type Service struct {
	Id              string         `json:"id,omitempty"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	NetworkItemType string         `json:"networkItemType"`
	NetworkItemId   string         `json:"networkItemId"`
	Type            string         `json:"type"`
	Routes          []*Route       `json:"routes"`
	Config          *ServiceConfig `json:"config"`
}

func (c *Client) CreateService(service *Service) (*Service, error) {
	serviceJson, err := json.Marshal(service)
	if err != nil {
		return nil, err
	}

	params := networkUrlParams(service.NetworkItemType, service.NetworkItemId)
	endpoint := fmt.Sprintf("%s/api/beta/services?%s", c.BaseURL, params.Encode())

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(serviceJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var s Service
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) GetService(serviceId, networkItemType, networkItemId string) (*Service, error) {
	params := networkUrlParams(networkItemType, networkItemId)
	params.Add("serviceId", serviceId)

	endpoint := fmt.Sprintf("%s/api/beta/services/single?%s", c.BaseURL, params.Encode())
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var s Service
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *Client) DeleteService(serviceId, networkItemType, networkItemId string) error {
	params := networkUrlParams(networkItemType, networkItemId).Encode()
	endpoint := fmt.Sprintf("%s/api/beta/services/%s?%s", c.BaseURL, serviceId, params)
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

func (c *Client) UpdateService(id, networkItemType, networkItemId string, service *Service) (*Service, error) {
	serviceJson, err := json.Marshal(service)
	if err != nil {
		return nil, err
	}

	params := networkUrlParams(networkItemType, networkItemId).Encode()
	endpoint := fmt.Sprintf("%s/api/beta/services/%s?%s", c.BaseURL, id, params)

	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(serviceJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var s Service
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func networkUrlParams(networkItemType string, networkItemId string) url.Values {
	params := url.Values{}
	params.Add("networkItemId", networkItemId)
	params.Add("networkItemType", networkItemType)
	return params
}
