package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type UserGroup struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	VpnRegionIds   []string `json:"vpnRegionIds"`
	InternetAccess string   `json:"internetAccess"`
	MaxDevice      int      `json:"maxDevice"`
	SystemSubnets  []string `json:"systemSubnets"`
}

func (c *Client) GetUserGroup(name string) (*UserGroup, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/user-groups", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var userGroups []UserGroup
	err = json.Unmarshal(body, &userGroups)
	if err != nil {
		return nil, err
	}
	for _, ug := range userGroups {
		if ug.Name == name {
			return &ug, nil
		}
	}
	return nil, nil
}
