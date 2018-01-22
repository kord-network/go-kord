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

	"github.com/ethereum/go-ethereum/common"
	meta "github.com/meta-network/go-meta"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr}
}

func (c *Client) Apply(tx *meta.Tx) (common.Hash, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return common.Hash{}, err
	}
	req, err := http.NewRequest("POST", c.addr+"/tx", bytes.NewReader(data))
	if err != nil {
		return common.Hash{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return common.Hash{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return common.Hash{}, fmt.Errorf("unxpected HTTP response: %s: %s", res.Status, body)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(string(body)), nil
}
