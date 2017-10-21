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

package api

import (
	"errors"
	"net/http"

	"github.com/flynn/flynn/pkg/httpclient"
	"github.com/flynn/flynn/pkg/stream"
	"github.com/ipfs/go-cid"
)

type Client struct {
	*httpclient.Client
}

func NewClient(baseURL string) *Client {
	return &Client{&httpclient.Client{
		HTTP:        http.DefaultClient,
		URL:         baseURL,
		ErrNotFound: errors.New("stream: not found"),
	}}
}

func (c *Client) ReadStream(name string, outputCh interface{}) (stream.Stream, error) {
	return c.Stream("GET", "/"+name, nil, outputCh)
}

func (c *Client) WriteStream(name string, v interface{}) (*cid.Cid, error) {
	var id cid.Cid
	return &id, c.Post("/"+name, v, &id)
}
