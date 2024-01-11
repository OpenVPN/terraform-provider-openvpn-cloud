package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ConnectionStatus string

type Connector struct {
	Id               string           `json:"id,omitempty"`
	Name             string           `json:"name"`
	NetworkItemId    string           `json:"networkItemId"`
	NetworkItemType  string           `json:"networkItemType"`
	VpnRegionId      string           `json:"vpnRegionId"`
	IPv4Address      string           `json:"ipV4Address"`
	IPv6Address      string           `json:"ipV6Address"`
	Profile          string           `json:"profile"`
	ConnectionStatus ConnectionStatus `json:"connectionStatus"`
}

type ConnectorPageResponse struct {
	Content          []Connector `json:"content"`
	NumberOfElements int         `json:"numberOfElements"`
	Page             int         `json:"page"`
	Size             int         `json:"size"`
	Success          bool        `json:"success"`
	TotalElements    int         `json:"totalElements"`
	TotalPages       int         `json:"totalPages"`
}

const (
	NetworkItemTypeHost    = "HOST"
	NetworkItemTypeNetwork = "NETWORK"
)

const (
	ConnectionStatusOffline ConnectionStatus = "OFFLINE"
	ConnectionStatusOnline  ConnectionStatus = "ONLINE"
)

func (c *Client) GetConnectorsByPage(page int, pageSize int) (ConnectorPageResponse, error) {
	endpoint := fmt.Sprintf("%s/api/beta/connectors/page?page=%d&size=%d", c.BaseURL, page, pageSize)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return ConnectorPageResponse{}, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return ConnectorPageResponse{}, err
	}

	var response ConnectorPageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ConnectorPageResponse{}, err
	}
	return response, nil
}

func (c *Client) GetAllConnectors() ([]Connector, error) {
	var allConnectors []Connector
	page := 1
	pageSize := 10

	for {
		response, err := c.GetConnectorsByPage(page, pageSize)
		if err != nil {
			return nil, err
		}

		allConnectors = append(allConnectors, response.Content...)

		if page >= response.TotalPages {
			break
		}
		page++
	}
	return allConnectors, nil
}

func (c *Client) GetConnectorByName(name string) (*Connector, error) {
	connectors, err := c.GetAllConnectors()
	if err != nil {
		return nil, err
	}

	for _, connector := range connectors {
		if connector.Name == name {
			return &connector, nil
		}
	}
	return nil, nil
}

func (c *Client) GetConnectorById(connectorId string) (*Connector, error) {
	connectors, err := c.GetAllConnectors()
	if err != nil {
		return nil, err
	}

	for _, connector := range connectors {
		if connector.Id == connectorId {
			return &connector, nil
		}
	}
	return nil, nil
}

func (c *Client) GetConnectorsForNetwork(networkId string) ([]Connector, error) {
	connectors, err := c.GetAllConnectors()
	if err != nil {
		return nil, err
	}

	var networkConnectors []Connector
	for _, connector := range connectors {
		if connector.NetworkItemId == networkId {
			networkConnectors = append(networkConnectors, connector)
		}
	}
	return networkConnectors, nil
}

func (c *Client) GetConnectorProfile(id string) (string, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/connectors/%s/profile", c.BaseURL, id), nil)
	if err != nil {
		return "", err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) CreateConnector(connector Connector, networkItemId string) (*Connector, error) {
	connectorJson, err := json.Marshal(connector)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/connectors?networkItemId=%s&networkItemType=%s", c.BaseURL, networkItemId, connector.NetworkItemType), bytes.NewBuffer(connectorJson))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var conn Connector
	err = json.Unmarshal(body, &conn)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (c *Client) DeleteConnector(connectorId string, networkItemId string, networkItemType string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/connectors/%s?networkItemId=%s&networkItemType=%s", c.BaseURL, connectorId, networkItemId, networkItemType), nil)
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	return err
}
