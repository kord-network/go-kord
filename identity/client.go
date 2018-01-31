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

package identity

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/meta-network/go-meta/graphql"
)

type Client struct {
	*graphql.Client
}

func NewClient(url string) *Client {
	return &Client{graphql.NewClient(url)}
}

func (c *Client) CreateIdentity(identity *Identity) error {
	query := `
mutation CreateIdentity($input: IdentityInput!) {
  createIdentity(input: $input) {
    id
    username
    owner
    signature
  }
}
`
	variables := graphql.Variables{"input": &IdentityInput{
		Username:  identity.Username,
		Owner:     identity.Owner.Hex(),
		Signature: hexutil.Encode(identity.Signature),
	}}
	return c.Do(query, variables, nil)
}

func (c *Client) Identity(filter *IdentityFilter) ([]*Identity, error) {
	query := `
query GetIdentity($filter: IdentityFilter!) {
  identity(filter: $filter) {
    id
    username
    owner
    signature
  }
}
`
	variables := graphql.Variables{"filter": filter}
	var v struct {
		Identities []*Identity `json:"identity"`
	}
	return v.Identities, c.Do(query, variables, &v)
}

func (c *Client) CreateClaim(claim *Claim) error {
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
		Issuer:    claim.Issuer.String(),
		Subject:   claim.Subject.String(),
		Property:  claim.Property,
		Claim:     claim.Claim,
		Signature: hexutil.Encode(claim.Signature),
	}}
	return c.Do(query, variables, nil)
}

func (c *Client) Claim(filter *ClaimFilter) ([]*Claim, error) {
	query := `
query GetClaim($filter: ClaimFilter!) {
  claim(filter: $filter) {
    id
    issuer
    subject
    property
    claim
    signature
  }
}
`
	variables := graphql.Variables{"filter": filter}
	var v struct {
		Claims []*Claim `json:"claim"`
	}
	return v.Claims, c.Do(query, variables, &v)
}
