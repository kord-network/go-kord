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
	"encoding/json"

	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/voc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func init() {
	voc.RegisterPrefix("id:", "http://schema.meta-network.io/identity/")
}

type ID struct {
	common.Hash
}

func NewID(hash common.Hash) ID {
	return ID{Hash: hash}
}

func HexToID(s string) ID {
	return NewID(common.HexToHash(s))
}

type Identity struct {
	Username  string
	Owner     common.Address
	Signature []byte
}

func (i *Identity) ID() ID {
	return NewID(crypto.Keccak256Hash([]byte(i.Username)))
}

type identityJSON struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Owner     string `json:"owner"`
	Signature string `json:"signature"`
}

func (i *Identity) MarshalJSON() ([]byte, error) {
	return json.Marshal(&identityJSON{
		ID:        i.ID().String(),
		Username:  i.Username,
		Owner:     i.Owner.String(),
		Signature: hexutil.Encode(i.Signature),
	})
}

func (i *Identity) UnmarshalJSON(b []byte) error {
	var v identityJSON
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*i = Identity{
		Username:  v.Username,
		Owner:     common.HexToAddress(v.Owner),
		Signature: common.FromHex(v.Signature),
	}
	return nil
}

type identityQuad struct {
	rdfType struct{} `quad:"@type > id:Identity"`

	ID        quad.IRI `quad:"@id"`
	Username  string   `quad:"id:username"`
	Owner     string   `quad:"id:owner"`
	Signature string   `quad:"id:signature"`
}

func (i *Identity) Quad() *identityQuad {
	return &identityQuad{
		ID:        quad.IRI(i.ID().String()),
		Username:  i.Username,
		Owner:     i.Owner.String(),
		Signature: hexutil.Encode(i.Signature),
	}
}

func (i *identityQuad) Identity() *Identity {
	return &Identity{
		Username:  i.Username,
		Owner:     common.HexToAddress(i.Owner),
		Signature: common.FromHex(i.Signature),
	}
}

type IdentityFilter struct {
	ID       *string `json:"id"`
	Username *string `json:"username"`
	Owner    *string `json:"owner"`
}

type IdentityInput struct {
	Username  string `json:"username"`
	Owner     string `json:"owner"`
	Signature string `json:"signature"`
}

type Claim struct {
	Graph     string
	Issuer    ID
	Subject   ID
	Property  string
	Claim     string
	Signature []byte
}

func (c *Claim) ID() common.Hash {
	return crypto.Keccak256Hash(
		c.Issuer.Hash[:],
		c.Subject.Hash[:],
		[]byte(c.Property),
		[]byte(c.Claim),
	)
}

type claimJSON struct {
	ID        string `json:"id"`
	Graph     string `json:"graph"`
	Issuer    string `json:"issuer"`
	Subject   string `json:"subject"`
	Property  string `json:"property"`
	Claim     string `json:"claim"`
	Signature string `json:"signature"`
}

func (c *Claim) MarshalJSON() ([]byte, error) {
	return json.Marshal(&claimJSON{
		ID:        c.ID().String(),
		Graph:     c.Graph,
		Issuer:    c.Issuer.String(),
		Subject:   c.Subject.String(),
		Property:  c.Property,
		Claim:     c.Claim,
		Signature: hexutil.Encode(c.Signature),
	})
}

func (c *Claim) UnmarshalJSON(b []byte) error {
	var v claimJSON
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*c = Claim{
		Graph:     v.Graph,
		Issuer:    HexToID(v.Issuer),
		Subject:   HexToID(v.Subject),
		Property:  v.Property,
		Claim:     v.Claim,
		Signature: common.FromHex(v.Signature),
	}
	return nil
}

type claimQuad struct {
	rdfType struct{} `quad:"@type > id:Claim"`

	ID        quad.IRI `quad:"@id"`
	Issuer    quad.IRI `quad:"id:issuer"`
	Subject   quad.IRI `quad:"id:subject"`
	Property  string   `quad:"id:property"`
	Claim     string   `quad:"id:claim"`
	Signature string   `quad:"id:signature"`
}

func (c *Claim) Quad() *claimQuad {
	return &claimQuad{
		ID:        quad.IRI(c.ID().String()),
		Issuer:    quad.IRI(c.Issuer.String()),
		Subject:   quad.IRI(c.Subject.String()),
		Property:  c.Property,
		Claim:     c.Claim,
		Signature: hexutil.Encode(c.Signature),
	}
}

func (c *claimQuad) ToClaim(graph string) *Claim {
	return &Claim{
		Graph:     graph,
		Issuer:    HexToID(string(c.Issuer)),
		Subject:   HexToID(string(c.Subject)),
		Property:  c.Property,
		Claim:     c.Claim,
		Signature: common.FromHex(c.Signature),
	}
}

type ClaimFilter struct {
	Graph    string  `json:"graph"`
	Issuer   *string `json:"issuer"`
	Subject  *string `json:"subject"`
	Property *string `json:"property"`
	Claim    *string `json:"claim"`
}

type ClaimInput struct {
	Graph     string `json:"graph"`
	Issuer    string `json:"issuer"`
	Subject   string `json:"subject"`
	Property  string `json:"property"`
	Claim     string `json:"claim"`
	Signature string `json:"signature"`
}
