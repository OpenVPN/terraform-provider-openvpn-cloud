package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Route struct {
	Id            string `json:"id,omitempty"`
	Type          string `json:"type,omitempty"`
	Subnet        string `json:"subnet,omitempty"`
	Description   string `json:"description,omitempty"`
	NetworkItemId string `json:"networkItemId,omitempty"`
}

type RoutePageResponse struct {
	Success          bool    `json:"success"`
	Content          []Route `json:"content"`
	TotalElements    int     `json:"totalElements"`
	TotalPages       int     `json:"totalPages"`
	NumberOfElements int     `json:"numberOfElements"`
	Page             int     `json:"page"`
	Size             int     `json:"size"`
}

const (
	RouteTypeIPV4   = "IP_V4"
	RouteTypeIPV6   = "IP_V6"
	RouteTypeDomain = "DOMAIN"
)

func (c *Client) GetRoutesByPage(networkId string, page int, size int) (RoutePageResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/networks/%s/routes/page?page=%d&size=%d", c.BaseURL, networkId, page, size), nil)
	if err != nil {
		return RoutePageResponse{}, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return RoutePageResponse{}, err
	}

	var response RoutePageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return RoutePageResponse{}, err
	}
	return response, nil
}

func (c *Client) GetAllRoutes(networkId string) ([]Route, error) {
	var allRoutes []Route
	pageSize := 10
	page := 1

	for {
		response, err := c.GetRoutesByPage(networkId, page, pageSize)
		if err != nil {
			return nil, err
		}

		allRoutes = append(allRoutes, response.Content...)

		if page >= response.TotalPages {
			break
		}
		page++
	}
	return allRoutes, nil
}

func (c *Client) GetNetworkRoute(networkId string, routeId string) (*Route, error) {
	routes, err := c.GetAllRoutes(networkId)
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
	networks, err := c.GetAllNetworks()
	if err != nil {
		return nil, err
	}

	for _, n := range networks {
		routes, err := c.GetAllRoutes(n.Id)
		if err != nil {
			continue
		}
		for _, r := range routes {
			if r.Id == routeId {
				r.NetworkItemId = n.Id
				return &r, nil
			}
		}
	}
	return nil, nil
}

func (c *Client) CreateRoute(networkId string, route Route) (*Route, error) {
	type newRoute struct {
		Description string `json:"description"`
		Subnet      string `json:"subnet"`
	}
	routeToCreate := newRoute{
		Description: route.Description,
		Subnet:      route.Subnet,
	}
	routeJson, err := json.Marshal(routeToCreate)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api/beta/networks/%s/routes", c.BaseURL, networkId),
		bytes.NewBuffer(routeJson),
	)
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

	// The API does not return the route Value, so we set it manually.
	r.Subnet = routeToCreate.Subnet
	return &r, nil
}

func (c *Client) UpdateRoute(networkId string, route Route) error {
	routeJson, err := json.Marshal(route)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/beta/networks/%s/routes/%s", c.BaseURL, networkId, route.Id),
		bytes.NewBuffer(routeJson),
	)
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	return err
}

func (c *Client) DeleteRoute(networkId string, routeId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/networks/%s/routes/%s", c.BaseURL, networkId, routeId), nil)
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	return err
}
