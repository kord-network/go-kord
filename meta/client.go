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

package meta

import (
	"context"

	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	client *rpc.Client
}

func NewClient(url string) (*Client, error) {
	c, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Client{c}, nil
}

func (c *Client) CreateGraph(ctx context.Context, id string) (common.Hash, error) {
	var hash common.Hash
	return hash, c.client.CallContext(ctx, &hash, "meta_createGraph", id)
}

func (c *Client) CommitGraph(ctx context.Context, id string) (common.Hash, error) {
	var hash common.Hash
	return hash, c.client.CallContext(ctx, &hash, "meta_commitGraph", id)
}

func (c *Client) SetGraph(ctx context.Context, hash common.Hash, sig []byte) error {
	return c.client.CallContext(ctx, nil, "meta_setGraph", hash, sig)
}

func (c *Client) QuadStore(name string) graph.QuadStore {
	return &clientQuadStore{c.client, name}
}

type clientQuadStore struct {
	client *rpc.Client
	name   string
}

func (c *clientQuadStore) ApplyDeltas(in []graph.Delta, opts graph.IgnoreOpts) error {
	return c.client.Call(nil, "meta_applyDeltas", c.name, in, opts)
}

func (c *clientQuadStore) Quad(graph.Value) quad.Quad {
	panic("method not implemented")
}

func (c *clientQuadStore) QuadIterator(quad.Direction, graph.Value) graph.Iterator {
	panic("method not implemented")
}

func (c *clientQuadStore) NodesAllIterator() graph.Iterator {
	panic("method not implemented")
}

func (c *clientQuadStore) QuadsAllIterator() graph.Iterator {
	panic("method not implemented")
}

func (c *clientQuadStore) ValueOf(quad.Value) graph.Value {
	panic("method not implemented")
}

func (c *clientQuadStore) NameOf(graph.Value) quad.Value {
	panic("method not implemented")
}

func (c *clientQuadStore) Size() int64 {
	panic("method not implemented")
}

func (c *clientQuadStore) OptimizeIterator(it graph.Iterator) (graph.Iterator, bool) {
	panic("method not implemented")
}

func (c *clientQuadStore) Close() error {
	panic("method not implemented")
}

func (c *clientQuadStore) QuadDirection(id graph.Value, d quad.Direction) graph.Value {
	panic("method not implemented")
}
