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
	IPV4Addresses []string `json:"ipv4Addresses"`
	IPV6Addresses []string `json:"ipv6Addresses"`
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

func (c *Client) GetDnsRecord(recordId string) (*DnsRecord, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/dns-records", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var records []DnsRecord
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		if r.Id == recordId {
			return &r, nil
		}
	}
	return nil, nil
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
