// This file is part of the go-meta library.
//
// Copyright (C) 2018 JAAK MUSIC LTD
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// If you have any questions please contact yo@jaak.io

package graphql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/neelance/graphql-go/errors"
)

type Variables map[string]interface{}

type Request struct {
	Query     string    `json:"query,omitempty"`
	Variables Variables `json:"variables,omitempty"`
}

type Response struct {
	Data   json.RawMessage      `json:"data,omitempty"`
	Errors []*errors.QueryError `json:"errors,omitempty"`
}

type Client struct {
	url string
}

func NewClient(url string) *Client {
	return &Client{url}
}

func (c *Client) Do(query string, variables Variables, out interface{}) error {
	data, err := json.Marshal(&Request{
		Query:     query,
		Variables: variables,
	})
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("POST", c.url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpRes, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()
	if httpRes.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(httpRes.Body)
		return fmt.Errorf("graphql: unexpected HTTP response: %s: %s", httpRes.Status, body)
	}
	var res Response
	if err := json.NewDecoder(httpRes.Body).Decode(&res); err != nil {
		return fmt.Errorf("graphql: error decoding GraphQL response: %s", err)
	}
	if len(res.Errors) > 0 {
		return fmt.Errorf("graphql: unexpected errors in GraphQL response: %v", res.Errors)
	}
	if out != nil {
		return json.Unmarshal(res.Data, out)
	}
	return nil
}
