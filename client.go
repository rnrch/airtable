// Copyright 2021 rnrch
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package airtable

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/go-resty/resty/v2"
	"go.uber.org/ratelimit"
)

const baseURL = "https://api.airtable.com/v0"

type Record struct {
	ID          string      `json:"id,omitempty"`
	Fields      interface{} `json:"fields"`
	CreatedTime string      `json:"createdTime,omitempty"`
}

type Records struct {
	Records []Record `json:"records"`
	Offset  string   `json:"offset,omitempty"`
}

type Client struct {
	client      *resty.Client
	rateLimiter ratelimit.Limiter
	APIKey      string
	BaseID      string
}

func NewClient(apiKey string, baseID string) *Client {
	return &Client{
		client:      resty.New(),
		rateLimiter: ratelimit.New(5),
		APIKey:      apiKey,
		BaseID:      baseID,
	}
}

func (c *Client) ListRecords(table string, params map[string]string) (Records, error) {
	records := Records{}
	u, err := url.Parse(baseURL)
	if err != nil {
		return records, err
	}
	u.Path = path.Join(u.Path, c.BaseID, table)
	c.rateLimiter.Take()
	resp, err := c.client.R().
		SetAuthToken(c.APIKey).
		SetQueryParams(params).
		Get(u.String())
	if err != nil {
		return records, err
	}
	if resp.StatusCode() != http.StatusOK {
		return records, fmt.Errorf("status: %s \nresp: %s", resp.Status(), resp.Body())
	}
	body := resp.Body()
	err = json.Unmarshal(body, &records)
	return records, err
}

func (c *Client) GetRecord(table string, id string) (Record, error) {
	r := Record{}
	u, err := url.Parse(baseURL)
	if err != nil {
		return r, err
	}
	u.Path = path.Join(u.Path, c.BaseID, table, id)
	c.rateLimiter.Take()
	resp, err := c.client.R().
		SetAuthToken(c.APIKey).
		Get(u.String())
	if err != nil {
		return r, err
	}
	if resp.StatusCode() != http.StatusOK {
		return r, fmt.Errorf("status: %s \nresp: %s", resp.Status(), resp.Body())
	}
	body := resp.Body()
	err = json.Unmarshal(body, &r)
	return r, err
}

func (c *Client) CreateRecords(table string, records Records) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	by, err := json.Marshal(records)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, c.BaseID, table)
	c.rateLimiter.Take()
	resp, err := c.client.R().
		SetAuthToken(c.APIKey).
		SetBody(string(by)).
		SetHeader("Content-Type", "application/json").
		Post(u.String())
	if err != nil {
		return err
	}
	log.Println(resp.Request.Body)
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("status: %s \nresp: %s", resp.Status(), resp.Body())
	}
	return nil
}

func (c *Client) PatchRecords(table string, records Records) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	u.Path = path.Join(u.Path, c.BaseID, table)
	c.rateLimiter.Take()
	resp, err := c.client.R().
		SetAuthToken(c.APIKey).
		SetBody(records).
		SetHeader("Content-Type", "application/json").
		Patch(u.String())
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Printf("%v", resp.Request.Body)
		return fmt.Errorf("status: %s, resp: %s", resp.Status(), resp.Body())
	}
	return nil
}

func (c *Client) DeleteRecords(table string, ids []string) error {
	records := resliceByNum(ids, 10)
	for _, reqs := range records {
		u, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		u.Path = path.Join(u.Path, c.BaseID, table)
		params := url.Values{"records[]": reqs}
		c.rateLimiter.Take()
		resp, err := c.client.R().
			SetAuthToken(c.APIKey).
			SetQueryParamsFromValues(params).
			Delete(u.String())
		if err != nil {
			return err
		}
		if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("status: %s \nresp: %s", resp.Status(), resp.Body())
		}
	}
	return nil
}

func resliceByNum(s []string, num int) [][]string {
	res := [][]string{}
	if num <= 0 {
		return res
	}
	for i, element := range s {
		if i%num == 0 {
			res = append(res, []string{})
		}
		res[len(res)-1] = append(res[len(res)-1], element)
	}
	return res
}
