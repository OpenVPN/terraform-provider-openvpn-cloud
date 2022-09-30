package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type VpnRegion struct {
	Id         string `json:"id"`
	Continent  string `json:"continent"`
	Country    string `json:"country"`
	CountryISO string `json:"countryIso"`
	RegionName string `json:"regionName"`
}

func (c *Client) GetVpnRegion(regionId string) (*VpnRegion, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/regions", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var vpnRegions []VpnRegion
	err = json.Unmarshal(body, &vpnRegions)
	if err != nil {
		return nil, err
	}
	for _, r := range vpnRegions {
		if r.Id == regionId {
			return &r, nil
		}
	}
	return nil, nil
}
