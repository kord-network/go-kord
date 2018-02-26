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
	"encoding/json"

	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/voc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func init() {
	voc.RegisterPrefix("meta:", "http://schema.meta-network.io/")
}

type ID struct {
	common.Address
}

func NewID(addr common.Address) ID {
	return ID{Address: addr}
}

func HexToID(s string) ID {
	return NewID(common.HexToAddress(s))
}

type GraphInput struct {
	ID string `json:"id"`
}

type Claim struct {
	Issuer    ID
	Subject   ID
	Property  string
	Claim     string
	Signature []byte
}

func (c *Claim) ID() common.Hash {
	return crypto.Keccak256Hash(
		c.Issuer.Address[:],
		c.Subject.Address[:],
		[]byte(c.Property),
		[]byte(c.Claim),
	)
}

type claimJSON struct {
	ID        string `json:"id"`
	Issuer    string `json:"issuer"`
	Subject   string `json:"subject"`
	Property  string `json:"property"`
	Claim     string `json:"claim"`
	Signature string `json:"signature"`
}

func (c *Claim) MarshalJSON() ([]byte, error) {
	return json.Marshal(&claimJSON{
		ID:        c.ID().String(),
		Issuer:    c.Issuer.Hex(),
		Subject:   c.Subject.Hex(),
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
	Issuer    quad.IRI `quad:"meta:issuer"`
	Subject   quad.IRI `quad:"meta:subject"`
	Property  string   `quad:"meta:property"`
	Claim     string   `quad:"meta:claim"`
	Signature string   `quad:"meta:signature"`
}

func (c *Claim) Quad() *claimQuad {
	return &claimQuad{
		ID:        quad.IRI(c.ID().String()),
		Issuer:    quad.IRI(c.Issuer.Hex()),
		Subject:   quad.IRI(c.Subject.Hex()),
		Property:  c.Property,
		Claim:     c.Claim,
		Signature: hexutil.Encode(c.Signature),
	}
}

func (c *claimQuad) ToClaim() *Claim {
	return &Claim{
		Issuer:    HexToID(string(c.Issuer)),
		Subject:   HexToID(string(c.Subject)),
		Property:  c.Property,
		Claim:     c.Claim,
		Signature: common.FromHex(c.Signature),
	}
}

type ClaimFilter struct {
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
