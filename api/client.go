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

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cayleygraph/cayley/graph"
)

type Client struct {
	graph.QuadStore

	addr string
	name string
}

func NewClient(addr, name string) *Client {
	return &Client{
		addr: addr,
		name: name,
	}
}

func (c *Client) Create() error {
	res, err := http.Post(fmt.Sprintf("%s/%s", c.addr, c.name), "", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("unxpected HTTP response: %s: %s", res.Status, body)
	}
	return nil
}

type Request struct {
	In   []graph.Delta
	Opts graph.IgnoreOpts
}

func (c *Client) ApplyDeltas(in []graph.Delta, opts graph.IgnoreOpts) error {
	data, err := json.Marshal(&Request{in, opts})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/apply-deltas", c.addr, c.name), bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("unxpected HTTP response: %s: %s", res.Status, body)
	}
	return nil
}
