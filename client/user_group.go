package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserGroup struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	ConnectAuth    string   `json:"connectAuth"`
	VpnRegionIds   []string `json:"vpnRegionIds"`
	InternetAccess string   `json:"internetAccess"`
	MaxDevice      int      `json:"maxDevice"`
	SystemSubnets  []string `json:"systemSubnets"`
}

func (c *Client) GetUserGroups() ([]UserGroup, error) {
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
	return userGroups, nil
}
func (c *Client) GetUserGroupByName(name string) (*UserGroup, error) {
	userGroups, err := c.GetUserGroups()
	if err != nil {
		return nil, err
	}
	for _, ug := range userGroups {
		if ug.Name == name {
			return &ug, nil
		}
	}
	return nil, fmt.Errorf("group %s does not exist", name)
}

func (c *Client) GetUserGroupById(id string) (*UserGroup, error) {
	userGroups, err := c.GetUserGroups()
	if err != nil {
		return nil, err
	}
	for _, ug := range userGroups {
		if ug.Id == id {
			return &ug, nil
		}
	}
	return nil, fmt.Errorf("group %s does not exist", id)
}

func (c *Client) CreateUserGroup(userGroup *UserGroup) (*UserGroup, error) {
	userGroupJson, err := json.Marshal(userGroup)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/user-groups", c.BaseURL), bytes.NewBuffer(userGroupJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var ug UserGroup
	err = json.Unmarshal(body, &ug)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}

func (c *Client) DeleteUserGroup(id string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/user-groups/%s", c.BaseURL, id), nil)
	if err != nil {
		return err
	}
	_, err = c.DoRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateUserGroup(id string, userGroup *UserGroup) (*UserGroup, error) {
	userGroupJson, err := json.Marshal(userGroup)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/beta/user-groups/%s", c.BaseURL, id), bytes.NewBuffer(userGroupJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var ug UserGroup
	err = json.Unmarshal(body, &ug)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}
