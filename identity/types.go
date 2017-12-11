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
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type ID struct {
	common.Hash
}

func NewID(hash common.Hash) ID {
	return ID{Hash: hash}
}

func HexToID(s string) ID {
	return NewID(common.HexToHash(s))
}

// Identity structure
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

// Claim structure
type Claim struct {
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
	Issuer    string `json:"issuer"`
	Subject   string `json:"subject"`
	Property  string `json:"property"`
	Claim     string `json:"claim"`
	Signature string `json:"signature"`
}

func (c *Claim) MarshalJSON() ([]byte, error) {
	return json.Marshal(&claimJSON{
		ID:        c.ID().String(),
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
		Issuer:    HexToID(v.Issuer),
		Subject:   HexToID(v.Subject),
		Property:  v.Property,
		Claim:     v.Claim,
		Signature: common.FromHex(v.Signature),
	}
	return nil
}
