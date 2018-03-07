// This file is part of the go-kord library.
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
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/kord-network/go-kord/graphql"
)

type Client struct {
	*graphql.Client
}

func NewClient(url string) *Client {
	return &Client{graphql.NewClient(url)}
}

func (c *Client) CreateGraph(graph string) (common.Hash, error) {
	query := `
mutation CreateGraph($input: GraphInput!) {
  createGraph(input: $input) {
    id
  }
}
`
	variables := graphql.Variables{"input": &GraphInput{
		ID: graph,
	}}
	res, err := c.Do(query, variables, nil)
	if err != nil {
		return common.Hash{}, err
	}
	return swarmHash(res)
}

func (c *Client) SetGraph(graph string, hash common.Hash, sig []byte) error {
	query := `
mutation SetGraph($input: SetGraphInput!) {
  setGraph(input: $input) {
    id
  }
}
`
	variables := graphql.Variables{"input": &SetGraphInput{
		ID:        graph,
		Hash:      hash.Hex(),
		Signature: hexutil.Encode(sig),
	}}
	_, err := c.Do(query, variables, nil)
	return err
}

func (c *Client) CreateClaim(graph string, claim *Claim) (common.Hash, error) {
	query := `
mutation CreateClaim($input: ClaimInput!) {
  createClaim(input: $input) {
    id
    issuer
    subject
    property
    claim
    signature
  }
}
`
	variables := graphql.Variables{"input": &ClaimInput{
		Graph:     graph,
		Issuer:    claim.Issuer.Hex(),
		Subject:   claim.Subject.Hex(),
		Property:  claim.Property,
		Claim:     claim.Claim,
		Signature: hexutil.Encode(claim.Signature),
	}}
	res, err := c.Do(query, variables, nil)
	if err != nil {
		return common.Hash{}, err
	}
	return swarmHash(res)
}

func (c *Client) Claim(graph string, filter *ClaimFilter) ([]*Claim, error) {
	query := `
query GetClaim($id: String!, $filter: ClaimFilter!) {
  graph(id: $id) {
    claim(filter: $filter) {
      id
      issuer
      subject
      property
      claim
      signature
    }
  }
}
`
	variables := graphql.Variables{"id": graph, "filter": filter}
	var v struct {
		Graph struct {
			Claims []*Claim `json:"claim"`
		} `json:"graph"`
	}
	if _, err := c.Do(query, variables, &v); err != nil {
		return nil, err
	}
	return v.Graph.Claims, nil
}

func swarmHash(res *graphql.Response) (common.Hash, error) {
	extension, ok := res.Extensions["kord"]
	if !ok {
		return common.Hash{}, errors.New("missing kord extension in GraphQL response")
	}
	v, ok := extension.(map[string]interface{})
	if !ok {
		return common.Hash{}, fmt.Errorf("unexpected kord extension type: %T", extension)
	}
	h, ok := v["swarmHash"].(string)
	if !ok {
		return common.Hash{}, fmt.Errorf("unexpected swarmHash type: %T", v["swarmHash"])
	}
	return common.HexToHash(h), nil
}
