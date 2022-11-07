package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Route struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	Subnet        string `json:"subnet"`
	Domain        string `json:"domain"`
	Value         string `json:"value"`
	NetworkItemId string `json:"networkItemId"`
	Description   string `json:"description"`
}

const (
	RouteTypeIPV4   = "IP_V4"
	RouteTypeIPV6   = "IP_V6"
	RouteTypeDomain = "DOMAIN"
)

func (c *Client) CreateRoute(networkId string, route Route) (*Route, error) {
	routeJson, err := json.Marshal(route)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/networks/%s/routes", c.BaseURL, networkId), bytes.NewBuffer(routeJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var r Route
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) DeleteRoute(networkId string, routeId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/networks/%s/routes/%s", c.BaseURL, networkId, routeId), nil)
	if err != nil {
		return err
	}
	_, err = c.DoRequest(req)
	return err
}

func (c *Client) GetRoutes(networkId string) ([]Route, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/networks/%s/routes", c.BaseURL, networkId), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var routes []Route
	err = json.Unmarshal(body, &routes)
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (c *Client) GetNetworkRoute(networkId string, routeId string) (*Route, error) {
	routes, err := c.GetRoutes(networkId)
	if err != nil {
		return nil, err
	}
	for _, r := range routes {
		if r.Id == routeId {
			return &r, nil
		}
	}
	return nil, nil
}

func (c *Client) GetRouteById(routeId string) (*Route, error) {
	networks, err := c.GetNetworks()
	if err != nil {
		return nil, err
	}
	for _, n := range networks {
		r, err := c.GetNetworkRoute(n.Id, routeId)
		if err != nil {
			return nil, err
		}
		if r != nil {
			r.NetworkItemId = n.Id
			return r, nil
		}
	}
	return nil, nil
}

func (c *Client) UpdateRoute(networkId string, route Route) error {
	routeJson, err := json.Marshal(route)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/beta/networks/%s/routes/%s", c.BaseURL, networkId, route.Id), bytes.NewBuffer(routeJson))
	if err != nil {
		return err
	}
	_, err = c.DoRequest(req)
	return err
}
