// This file is part of the go-meta library.
//
// Copyright (C) 2017 JAAK MUSIC LTD
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

package identity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL}
}

func (c *Client) CreateIdentity(name string) (*Identity, error) {
	form := make(url.Values)
	form.Set("name", name)
	res, err := http.PostForm(c.baseURL, form)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected HTTP response: %s: %s", res.Status, body)
	}
	var identity Identity
	if err := json.NewDecoder(res.Body).Decode(&identity); err != nil {
		return nil, err
	}
	return &identity, nil
}

func (c *Client) GetIdentity(name string) (*Identity, error) {
	res, err := http.Get(c.baseURL + "/" + name)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected HTTP response: %s: %s", res.Status, body)
	}
	var identity Identity
	if err := json.NewDecoder(res.Body).Decode(&identity); err != nil {
		return nil, err
	}
	return &identity, nil
}

func (c *Client) UpdateIdentity(name string, aux map[string]string) (*Identity, error) {
	data, err := json.Marshal(aux)
	if err != nil {
		return nil, err
	}
	res, err := http.Post(c.baseURL+"/"+name, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("unexpected HTTP response: %s: %s", res.Status, body)
	}
	var identity Identity
	if err := json.NewDecoder(res.Body).Decode(&identity); err != nil {
		return nil, err
	}
	return &identity, nil
}
