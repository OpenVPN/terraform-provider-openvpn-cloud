package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type DnsRecord struct {
	Id            string   `json:"id"`
	Domain        string   `json:"domain"`
	Description   string   `json:"description"`
	IPV4Addresses []string `json:"ipv4Addresses"`
	IPV6Addresses []string `json:"ipv6Addresses"`
}

type DnsRecordPageResponse struct {
	Content          []DnsRecord `json:"content"`
	NumberOfElements int         `json:"numberOfElements"`
	Page             int         `json:"page"`
	Size             int         `json:"size"`
	Success          bool        `json:"success"`
	TotalElements    int         `json:"totalElements"`
	TotalPages       int         `json:"totalPages"`
}

func (c *Client) GetDnsRecordsByPage(page int, pageSize int) (DnsRecordPageResponse, error) {
	endpoint := fmt.Sprintf("%s/api/beta/dns-records/page?page=%d&size=%d", c.BaseURL, page, pageSize)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return DnsRecordPageResponse{}, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return DnsRecordPageResponse{}, err
	}

	var response DnsRecordPageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return DnsRecordPageResponse{}, err
	}
	return response, nil
}

func (c *Client) GetDnsRecord(recordId string) (*DnsRecord, error) {
	pageSize := 10
	page := 1

	for {
		response, err := c.GetDnsRecordsByPage(page, pageSize)
		if err != nil {
			return nil, err
		}

		for _, record := range response.Content {
			if record.Id == recordId {
				return &record, nil
			}
		}

		if page >= response.TotalPages {
			break
		}
		page++
	}
	return nil, fmt.Errorf("DNS record with ID %s not found", recordId)
}

func (c *Client) CreateDnsRecord(record DnsRecord) (*DnsRecord, error) {
	recordJson, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/dns-records", c.BaseURL), bytes.NewBuffer(recordJson))
	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}

	var d DnsRecord
	err = json.Unmarshal(body, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (c *Client) UpdateDnsRecord(record DnsRecord) error {
	recordJson, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/beta/dns-records/%s", c.BaseURL, record.Id), bytes.NewBuffer(recordJson))
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	return err
}

func (c *Client) DeleteDnsRecord(recordId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/dns-records/%s", c.BaseURL, recordId), nil)
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	return err
}
