package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id        string   `json:"id"`
	Username  string   `json:"username"`
	Role      string   `json:"role"`
	Email     string   `json:"email"`
	AuthType  string   `json:"authType"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	GroupId   string   `json:"groupId"`
	Status    string   `json:"status"`
	Devices   []Device `json:"devices"`
}

type Device struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IPv4Address string `json:"ipV4Address"`
	IPv6Address string `json:"ipV6Address"`
}

func (c *Client) CreateUser(user User) (*User, error) {
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/users", c.BaseURL), bytes.NewBuffer(userJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var u User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) GetUser(username string, role string) (*User, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/users", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var users []User
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}
	for _, u := range users {
		if u.Username == username && u.Role == role {
			return &u, nil
		}
	}
	return nil, nil
}

func (c *Client) GetUserById(userId string) (*User, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/users/%s", c.BaseURL, userId), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var u User
	err = json.Unmarshal(body, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *Client) DeleteUser(userId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/users/%s", c.BaseURL, userId), nil)
	if err != nil {
		return err
	}
	_, err = c.DoRequest(req)
	return err
}
