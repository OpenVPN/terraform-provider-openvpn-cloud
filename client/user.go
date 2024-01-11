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

type UserPageResponse struct {
	Content          []User `json:"content"`
	NumberOfElements int    `json:"numberOfElements"`
	Page             int    `json:"page"`
	Size             int    `json:"size"`
	Success          bool   `json:"success"`
	TotalElements    int    `json:"totalElements"`
	TotalPages       int    `json:"totalPages"`
}

type Device struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IPv4Address string `json:"ipV4Address"`
	IPv6Address string `json:"ipV6Address"`
}

func (c *Client) GetUsersByPage(page int, pageSize int) (UserPageResponse, error) {
	endpoint := fmt.Sprintf("%s/api/beta/users/page?page=%d&size=%d", c.BaseURL, page, pageSize)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return UserPageResponse{}, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return UserPageResponse{}, err
	}

	var response UserPageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return UserPageResponse{}, err
	}
	return response, nil
}

func (c *Client) GetUser(username string, role string) (*User, error) {
	pageSize := 10
	page := 1

	for {
		response, err := c.GetUsersByPage(page, pageSize)
		if err != nil {
			return nil, err
		}

		for _, user := range response.Content {
			if user.Username == username && user.Role == role {
				return &user, nil
			}
		}

		if page >= response.TotalPages {
			break
		}
		page++
	}
	return nil, fmt.Errorf("user with username %s and role %s not found", username, role)
}

func (c *Client) GetUserById(userId string) (*User, error) {
	pageSize := 10
	page := 1

	for {
		response, err := c.GetUsersByPage(page, pageSize)
		if err != nil {
			return nil, err
		}

		for _, user := range response.Content {
			if user.Id == userId {
				return &user, nil
			}
		}

		if page >= response.TotalPages {
			break
		}
		page++
	}
	return nil, fmt.Errorf("user with ID %s not found", userId)
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

func (c *Client) UpdateUser(user User) error {
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/beta/users/%s", c.BaseURL, user.Id), bytes.NewBuffer(userJson))
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteUser(userId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/users/%s", c.BaseURL, userId), nil)
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	return err
}
